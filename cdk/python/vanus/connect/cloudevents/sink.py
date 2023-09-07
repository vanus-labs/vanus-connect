# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import httpx
from cloudevents.abstract import CloudEvent
from cloudevents.conversion import to_structured

from ..cdk.sink import Sink


class CloudEventSink(Sink):
    def __init__(self, endpoint: str):
        super().__init__()
        self._endpoint = endpoint
        # self._client = httpx.AsyncClient(http2=True)

    async def on_event(self, event: CloudEvent):
        # Creates the HTTP request representation of the CloudEvent in binary content mode
        headers, body = to_structured(event)

        # POST
        async with httpx.AsyncClient(http2=True, timeout=httpx.Timeout(10.0)) as client:
            resp = await client.post(self._endpoint, content=body, headers=headers)

        # TODO: reuse http client
        # resp = await self._client.post(self.sink_endpoint, content=body, headers=headers)

        if resp.status_code / 100 != 2:
            raise RuntimeError(f"Failed to send event: {resp.status_code} {resp.text}")

        # TODO: return CloudEvent if present.
        # return resp
