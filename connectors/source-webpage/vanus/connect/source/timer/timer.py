# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import asyncio
import copy
import logging
from typing import Any, Dict

from cloudevents.http import CloudEvent
from vanus.connect.cdk import Source

attrs_proto = {
    "source": "https://github.com/vanus-labs/vanus-connect/connectors/source-webpage",
    "type": "ai.vanus.connect.source.timer",
}


class TimerSource(Source):
    def __init__(self, interval: float, data: Dict[str, Any]):
        super().__init__()
        self._interval = interval
        self._data = data

    async def start(self):
        """start the source"""
        await super().start()

        asyncio.create_task(self.run_timer(self._interval))

    async def run_timer(self, interval: float):
        try:
            while True:
                for i in range(20):
                    if i != 0:
                        # TODO: backoff
                        await asyncio.sleep(30)

                    try:
                        await self.on_time()
                        break
                    except Exception as e:
                        logging.warn("Sending event failed.", exc_info=e)

                await asyncio.sleep(interval)
        except asyncio.CancelledError:
            logging.info("timer cancelled")

    async def on_time(self):
        event = CloudEvent.create(copy.copy(attrs_proto), copy.deepcopy(self._data))
        await self.on_event(event)
