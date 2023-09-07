# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from abc import abstractmethod
from typing import TypeVar

from cloudevents.abstract import CloudEvent

from .sink import Sink
from .source import Source

AnyMiddleware = TypeVar("AnyMiddleware", bound="Middleware")


class Middleware(Source, Sink):
    def __init__(self):
        super().__init__()

    async def start(self):
        """start the sink"""
        await Source.start(self)

    async def on_event(self, event: CloudEvent):
        """receive event from source"""
        new_event = await self._on_event(event)
        return await Source.on_event(self, new_event)

    @abstractmethod
    async def _on_event(self, event: CloudEvent) -> CloudEvent:
        """receive event from source"""


class ServiceCallingMiddleware(Middleware):
    def __init__(self, service: Sink):
        super().__init__()
        self._service = service

    async def _on_event(self, event: CloudEvent) -> CloudEvent:
        resp = await self._service.on_event(event)
        assert resp is not None
        return resp
