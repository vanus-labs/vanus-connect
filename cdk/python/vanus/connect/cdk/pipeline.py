# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from typing import Generic, TypeVar, Union

from .middleware import ServiceCallingMiddleware
from .sink import AnySink, Sink
from .source import AnySource, Source

T = TypeVar("T", bound=Union[Source, Sink])


class Pipeline(Generic[AnySource, T]):
    def __init__(self, source: AnySource, last: T):
        self._source = source
        self._last = last

    def then(self, sink: AnySink) -> "Pipeline[AnySource, AnySink]":
        assert isinstance(self._last, Source)
        return Pipeline(self._source, self._last.connect(sink))

    def call(self, service: AnySink) -> "Pipeline[AnySource, AnySink]":
        assert isinstance(self._last, Source)
        return Pipeline(self._source, self._last.connect(ServiceCallingMiddleware(service)))

    async def start(self):
        await self._source.start()


def build_pipeline(source: AnySource) -> Pipeline[AnySource, AnySource]:
    return Pipeline(source, source)
