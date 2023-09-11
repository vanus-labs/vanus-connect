# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import asyncio
from http import HTTPStatus
from typing import Any, Awaitable, Callable

from cloudevents.abstract import CloudEvent
from cloudevents.conversion import to_structured
from cloudevents.http import from_http
from hypercorn.asyncio import serve
from hypercorn.config import Config
from quart import Quart, ResponseReturnValue, request
from quart.views import View

from ..cdk.source import Source

EventListener = Callable[[CloudEvent], Awaitable[Any]]


class CloudEventHandler(View):
    methods = ["POST"]

    # TODO: reuse handler after bug fixing of quart(>0.18.4)
    # init_every_request = False

    def __init__(self, on_event: EventListener) -> None:
        super().__init__()
        self._on_event = on_event

    async def dispatch_request(self, **kwargs: Any) -> ResponseReturnValue:
        # Create a CloudEvent
        event = from_http(request.headers, await request.get_data())

        # Process the CloudEvent
        new_event = await self.on_event(event)

        if new_event is None:
            return "", HTTPStatus.NO_CONTENT

        headers, body = to_structured(new_event)
        return body, HTTPStatus.OK, headers

    async def on_event(self, event):
        return await self._on_event(event)


class CloudEventSource(Source):
    def __init__(self, port: int):
        super().__init__()
        self._port = port

    async def start(self):
        try:
            import uvloop

            uvloop.install()
        except ImportError:
            pass

        config = Config()
        config.bind = [f"0.0.0.0:{self._port}"]

        app = Quart(__name__)
        app.add_url_rule("/", view_func=CloudEventHandler.as_view("source", self.on_event))

        asyncio.create_task(serve(app, config))
