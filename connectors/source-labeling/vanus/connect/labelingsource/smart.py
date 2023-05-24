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


import httpx
from cloudevents.http import CloudEvent
from jsonpath_ng import parse
from vanus.connect.customsource import Message

from .labeling import LabelMaker


class HttpSmartLabelMaker:
    def __init__(self, source_path, target_path, api_endpoint, ce_source=None, ce_type=None, *args, **kwargs):
        if ce_source is None:
            ce_source = "https://github.com/vanus-labs/vanus-connect/connectors/source-labeling"
        if ce_type is None:
            ce_type = "ai.vanus.connect.source.labeling"

        self._source_expr = parse(source_path)
        self._target_expr = parse(target_path)
        self._api_endpoint = api_endpoint
        self._label_maker = LabelMaker(*args, **kwargs)
        self._ce_source = ce_source
        self._ce_type = ce_type

    async def alabel(self, msg: Message) -> CloudEvent:
        for match in self._source_expr.find(msg):
            async with httpx.AsyncClient(http2=True, timeout=httpx.Timeout(120.0)) as client:
                resp = await client.post(self._api_endpoint, json={"content": match.value})

            if resp.is_success:
                body = resp.json()
                if body.get("status") == "normal":
                    self._target_expr.update_or_create(msg, body["labels"])
                    break

            labels = self._label_maker.label(match.value)
            self._target_expr.update_or_create(msg, labels)
            break

        attributes = {
            "source": self._ce_source,
            "type": self._ce_type,
        }
        event = CloudEvent.create(attributes, msg)

        return event
