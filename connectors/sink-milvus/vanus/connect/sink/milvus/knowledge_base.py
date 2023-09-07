# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import logging
from typing import cast

import openai

from .asynctool import ensure_async
from .milvus import Milvus, MilvusVectorStore

logger = logging.getLogger(__name__)


class KnowledgeBase:
    def __init__(self, **kwargs):
        openai.api_key = kwargs.pop("api_key")

        name = kwargs["vector_store"].pop("name")
        if name != "zilliz_ada":
            raise ValueError(f"unknow vector_store name: {name}")
        self._embedding_key = kwargs.pop("embedding_key")
        self.model = kwargs.pop("model")
        self._vector_store = _new_vector_store(**kwargs["vector_store"])

    async def load_document_from_texts(self, data):
        if self._embedding_key in data:
            emb = data[self._embedding_key]
            try:
                embeddings = await get_embedding(emb, self.model)
                data[self._embedding_key] = embeddings
            except Exception as e:
                raise e

            await ensure_async(self._vector_store.insert_chunk)(data)

            return True, "Success(embeddings)"
        await ensure_async(self._vector_store.insert_chunk)(data)
        return True, "Success(non-embeddings)"


def _new_vector_store(type, **kwargs):
    if type == "milvus":
        db = _new_milvus(**kwargs)
        return MilvusVectorStore(cast(Milvus, db))
    else:
        raise ValueError(f"unknown vector store type: {type}")


def _new_milvus(collection, **kwargs):
    return Milvus(collection_name=collection, connection_args=kwargs)


async def get_embedding(text, model="text-embedding-ada-002") -> float:
    return openai.Embedding.create(input=[text], model=model)["data"][0]["embedding"]  # type: ignore
