# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from abc import ABC, abstractmethod
from typing import Optional, TypeVar

from cloudevents.abstract import CloudEvent

from .sink import AnySink, Sink

AnySource = TypeVar("AnySource", bound="Source")


class Source(ABC):
    def __init__(self):
        super().__init__()
        self._sink: Optional[Sink] = None

    @abstractmethod
    async def start(self):
        """start the source"""
        if self._sink is None:
            raise RuntimeError("sink not connected")
        await self._sink.start()

    def connect(self, sink: AnySink) -> AnySink:
        """connect to sink"""
        self._sink = sink
        return sink

    async def on_event(self, event: CloudEvent):
        assert self._sink is not None
        return await self._sink.on_event(event)
