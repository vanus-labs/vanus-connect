# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from typing import Any, Dict, List, cast

import openai
from cloudevents.abstract import CloudEvent as AnyCloudEvent
from vanus.connect.cdk import Sink


class OpenAIEmbeddingService(Sink):
    def __init__(self, text_key: str, vector_key: str, **kwargs: Any) -> None:
        self._text_key = text_key
        self._vector_key = vector_key
        self._embedding_args = kwargs

    async def on_event(self, event: AnyCloudEvent) -> AnyCloudEvent:
        data = cast(Dict[str, Any], event.get_data())
        data[self._vector_key] = await self.embed_text(data[self._text_key])
        return event

    async def embed_text(self, text: str) -> List[float]:
        result = await openai.Embedding.acreate(input=[text], **self._embedding_args)
        return result["data"][0]["embedding"]  # type: ignore
