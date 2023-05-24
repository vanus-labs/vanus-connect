# Copyright 2023 Linkall Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from http import HTTPStatus
from typing import Any, Awaitable, Callable, Dict, Optional

import httpx
from cloudevents.abstract import AnyCloudEvent
from cloudevents.conversion import to_structured
from cloudevents.http import from_http
from quart import Quart, request
from quart.views import View


class Source:
    def __init__(self, sink_endpoint: str):
        self.sink_endpoint = sink_endpoint
        # self._client = httpx.AsyncClient(http2=True)

    async def on_event(self, event):
        # Creates the HTTP request representation of the CloudEvent in binary content mode
        headers, body = to_structured(event)

        # POST
        async with httpx.AsyncClient(http2=True, timeout=httpx.Timeout(10.0)) as client:
            resp = await client.post(self.sink_endpoint, content=body, headers=headers)

        # TODO: reuse http client
        # resp = await self._client.post(self.sink_endpoint, content=body, headers=headers)

        if resp.status_code / 100 != 2:
            raise RuntimeError(f"Failed to send event: {resp.status_code} {resp.text}")

        # TODO: return
        return resp


EventListener = Callable[[AnyCloudEvent], Awaitable[Any]]


class CloudEventHandler(View):
    methods = ["POST"]

    # TODO: reuse handler after bug fixing of quart(>0.18.4)
    # init_every_request = False

    def __init__(self, on_event: EventListener) -> None:
        super().__init__()
        self._on_event = on_event

    async def dispatch_request(self, **kwargs):
        # Create a CloudEvent
        event = from_http(request.headers, await request.get_data())

        # Process the CloudEvent
        await self.on_event(event)

        return "", HTTPStatus.NO_CONTENT

    async def on_event(self, event):
        return await self._on_event(event)


class CloudEventsSource(Source):
    def __init__(self, sink_endpoint: str, name=None, app=None):
        super().__init__(sink_endpoint)

        if app is None:
            if name is None:
                name = __name__
            app = Quart(name)

        self.app: Quart = app

        self.register_event_handler(self.on_event)

    def register_event_handler(self, on_event):
        view_func = CloudEventHandler.as_view("source", on_event)
        self.app.add_url_rule("/", view_func=view_func)
        self.app.add_url_rule("/<path:path>", view_func=view_func)


SyncEventHandler = Callable[[AnyCloudEvent], Optional[AnyCloudEvent]]
AsyncEventHandler = Callable[[AnyCloudEvent], Awaitable[Optional[AnyCloudEvent]]]


class CustomEventHandler(CloudEventHandler):
    def __init__(
        self,
        on_event: EventListener,
        async_handler: Optional[AsyncEventHandler] = None,
        sync_handler: Optional[SyncEventHandler] = None,
    ) -> None:
        super().__init__(on_event)
        if async_handler is not None:
            self._handle = async_handler
        elif sync_handler is not None:

            async def handler_impl(event: AnyCloudEvent) -> Optional[AnyCloudEvent]:
                return sync_handler(event)

            self._handle = handler_impl
        else:
            raise RuntimeError("no handler")

    async def on_event(self, event):
        # Handle CloudEvent and return a new CloudEvent
        new_event = await self._handle(event)

        if new_event is None:
            new_event = event

        return await super().on_event(new_event)


class CustomSource(CloudEventsSource):
    def __init__(
        self,
        sink_endpoint: str,
        async_handler: Optional[AsyncEventHandler] = None,
        sync_handler: Optional[SyncEventHandler] = None,
        **kwargs,
    ):
        self._handler_kwargs = {"async_handler": async_handler, "sync_handler": sync_handler}
        super().__init__(sink_endpoint, **kwargs)

    def register_event_handler(self, on_event: EventListener):
        view_func = CustomEventHandler.as_view("source", on_event, **self._handler_kwargs)
        self.app.add_url_rule("/", view_func=view_func)
        self.app.add_url_rule("/<path:path>", view_func=view_func)


Message = Dict[str, Any]
SyncMessageHandler = Callable[[Message], AnyCloudEvent]
AsyncMessageHandler = Callable[[Message], Awaitable[AnyCloudEvent]]


class CustomHTTPHandler(CloudEventHandler):
    def __init__(
        self,
        on_event: EventListener,
        async_handler: Optional[AsyncMessageHandler] = None,
        sync_handler: Optional[SyncMessageHandler] = None,
    ) -> None:
        super().__init__(on_event)
        if async_handler is not None:
            self._handle = async_handler
        elif sync_handler is not None:

            async def handler_impl(msg: Message) -> AnyCloudEvent:
                return sync_handler(msg)

            self._handle = handler_impl
        else:
            raise RuntimeError("no handler")

    async def dispatch_request(self, **kwargs):
        msg = await request.get_json()

        # Process the message
        await self.on_message(msg)

        return "", HTTPStatus.NO_CONTENT

    async def on_message(self, msg: Message):
        event = await self._handle(msg)
        return await super().on_event(event)


class CustomHTTPSource(CloudEventsSource):
    def __init__(
        self,
        sink_endpoint: str,
        async_handler: Optional[AsyncMessageHandler] = None,
        sync_handler: Optional[SyncMessageHandler] = None,
        **kwargs,
    ):
        self._handler_kwargs = {"async_handler": async_handler, "sync_handler": sync_handler}
        super().__init__(sink_endpoint, **kwargs)

    def register_event_handler(self, on_event):
        view_func = CustomHTTPHandler.as_view("source", on_event, **self._handler_kwargs)
        self.app.add_url_rule("/", view_func=view_func)
        self.app.add_url_rule("/<path:path>", view_func=view_func)
