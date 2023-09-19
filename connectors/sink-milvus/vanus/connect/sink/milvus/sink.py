# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import asyncio
import functools
import logging
from typing import Any, Dict, Union, cast
from uuid import uuid4

from cloudevents.abstract import CloudEvent as AnyCloudEvent
from pymilvus import Collection, MilvusException, connections
from vanus.connect.cdk import Sink

logger = logging.getLogger(__name__)


def _validate_mode(mode: str):
    if mode not in {"insert", "upsert"}:
        raise ValueError(f"invalid mode: {mode}")


class MilvusSink(Sink):
    def __init__(self, collection: str, mode: str = "insert", load_args: Dict[str, Any] = {}, **kwargs: Any) -> None:
        super().__init__()

        _validate_mode(mode)

        self._collection_name = collection
        self._mode = mode
        self._connection_args = kwargs
        self._load_args = load_args

    async def start(self):
        self._alias = self._create_connection_alias(self._connection_args)
        self._col = Collection(self._collection_name, using=self._alias)

        if self._mode == "upsert":
            self._load_func = self._col.upsert
        else:  # self._mode == "insert"
            self._load_func = self._col.insert

    async def on_event(self, event: AnyCloudEvent) -> None:
        assert hasattr(self, "_load_func")

        data = cast(Dict[str, Any], event.get_data())

        await asyncio.get_running_loop().run_in_executor(
            None, functools.partial(self._load_func, data, **self._load_args)
        )

    def _create_connection_alias(self, connection_args: dict) -> str:
        """Create the connection to the Milvus server."""

        # Grab the connection arguments that are used for checking existing connection
        host: str = connection_args.get("host", None)
        port: Union[str, int] = connection_args.get("port", None)
        address: str = connection_args.get("address", None)
        uri: str = connection_args.get("uri", None)
        user = connection_args.get("user", None)

        # Order of use is host/port, uri, address
        if host is not None and port is not None:
            given_address = str(host) + ":" + str(port)
        elif uri is not None:
            given_address = uri.split("https://")[1]
        elif address is not None:
            given_address = address
        else:
            given_address = None
            logger.debug("Missing standard address type for reuse attempt")

        # User defaults to empty string when getting connection info
        if user is None:
            user = ""

        # If a valid address was given, then check if a connection exists
        if given_address is not None:
            for con in connections.list_connections():
                addr = connections.get_connection_addr(con[0])
                if (
                    con[1]
                    and ("address" in addr)
                    and (addr["address"] == given_address)
                    and ("user" in addr)
                    and (addr["user"] == user)
                ):
                    logger.debug("Using previous connection: %s", con[0])
                    return con[0]

        # Generate a new connection if one doesn't exist
        alias = uuid4().hex
        try:
            connections.connect(alias=alias, **connection_args)
            logger.debug("Created new connection using: %s", alias)
            return alias
        except MilvusException as e:
            logger.error("Failed to create new connection using: %s", alias)
            raise e
