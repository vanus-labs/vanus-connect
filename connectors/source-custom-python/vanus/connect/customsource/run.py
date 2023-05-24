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

from typing import Optional

from hypercorn.config import Config

from .source import (
    AsyncEventHandler,
    AsyncMessageHandler,
    CustomHTTPSource,
    CustomSource,
    SyncEventHandler,
    SyncMessageHandler,
)


def _run(app, config):
    import asyncio

    from hypercorn.asyncio import serve

    try:
        import uvloop

        uvloop.install()
    except ImportError:
        pass

    asyncio.run(serve(app, config))


def _run_source(port, source):
    config = Config()
    config.bind = [f"0.0.0.0:{port}"]
    _run(source.app, config)


def run_source(
    port,
    sink_endpoint,
    async_handler: Optional[AsyncEventHandler] = None,
    sync_handler: Optional[SyncEventHandler] = None,
    name=None,
):
    source = CustomSource(sink_endpoint, async_handler=async_handler, sync_handler=sync_handler, name=name)
    _run_source(port, source)


def run_http_source(
    port,
    sink_endpoint,
    async_handler: Optional[AsyncMessageHandler] = None,
    sync_handler: Optional[SyncMessageHandler] = None,
    name=None,
):
    source = CustomHTTPSource(sink_endpoint, async_handler=async_handler, sync_handler=sync_handler, name=name)
    _run_source(port, source)
