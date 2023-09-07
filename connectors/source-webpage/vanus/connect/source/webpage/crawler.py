# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import copy

import bs2json
import httpx
from bs4 import BeautifulSoup
from cloudevents.abstract import CloudEvent as AnyCloudEvent
from cloudevents.http import CloudEvent
from vanus.connect.cdk import Sink

attrs_proto = {
    "source": "https://github.com/vanus-labs/vanus-connect/connectors/source-webpage",
    "type": "ai.vanus.connect.source.webpage",
}


class CrawlerService(Sink):
    async def on_event(self, event: AnyCloudEvent) -> CloudEvent:
        data = event.get_data()
        assert data is not None

        result = await self.crawl(data["url"])

        return CloudEvent.create(copy.copy(attrs_proto), result)

    async def crawl(self, url: str):
        async with httpx.AsyncClient() as client:
            response = await client.get(url)

        if response.status_code != 200:
            raise Exception(f"Failed to crawl {url}")

        soup = BeautifulSoup(response.content, "lxml")
        return bs2json.to_json(soup)
