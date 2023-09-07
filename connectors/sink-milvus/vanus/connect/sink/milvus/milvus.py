# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from __future__ import annotations

import logging
from typing import Any, Optional, Union
from uuid import uuid4

from pymilvus import Collection, utility

logger = logging.getLogger(__name__)

DEFAULT_MILVUS_CONNECTION = {
    "host": "localhost",
    "port": "19530",
    "user": "",
    "password": "",
    "secure": False,
}


class Milvus:
    """Initialize wrapper around the milvus vector database."""

    def __init__(
        self,
        collection_name: str = "LangChainCollection",
        connection_args: Optional[dict[str, Any]] = None,
        consistency_level: str = "Session",
        drop_old: Optional[bool] = False,
    ):
        """Initialize the Milvus vector store."""

        # Default search params when one is not provided.

        self.collection_name = collection_name
        self.consistency_level = consistency_level

        # In order for a collection to be compatible, pk needs to be auto'id and int
        # In order for compatibility, the text field will need to be called "text"
        # In order for compatibility, the vector field needs to be called "vector"
        self.fields: list[str] = []
        # Create the connection to the server
        if connection_args is None:
            connection_args = DEFAULT_MILVUS_CONNECTION
        self.alias = self._create_connection_alias(connection_args)
        self.col: Optional[Collection] = None

        # Grab the existing collection if it exists
        try:
            if utility.has_collection(self.collection_name, using=self.alias):
                self.col = Collection(
                    self.collection_name,
                    using=self.alias,
                )
        except Exception:
            raise UserWarning("not found existing collection")
        # If need to drop old, drop it
        if drop_old and isinstance(self.col, Collection):
            self.col.drop()
            self.col = None

    def _create_connection_alias(self, connection_args: dict) -> str:
        """Create the connection to the Milvus server."""
        from pymilvus import MilvusException, connections

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
            logger.debug("Missing standard address type for reuse atttempt")

        # User defaults to empty string when getting connection info
        if user is not None:
            tmp_user = user
        else:
            tmp_user = ""

        # If a valid address was given, then check if a connection exists
        if given_address is not None:
            for con in connections.list_connections():
                addr = connections.get_connection_addr(con[0])
                if (
                    con[1]
                    and ("address" in addr)
                    and (addr["address"] == given_address)
                    and ("user" in addr)
                    and (addr["user"] == tmp_user)
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


class MilvusVectorStore:
    def __init__(self, vector_store: Milvus):
        self._vs = vector_store

    def insert_chunk(
        self,
        data,
        timeout: Optional[int] = None,
        **kwargs,
    ):
        assert isinstance(self._vs.col, Collection)

        # FIXME: pk is auto
        self._vs.col.insert(data, timeout=timeout, **kwargs)
