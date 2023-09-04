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


import json
from io import StringIO
from typing import Optional

from actrie import Matcher
from cloudevents.abstract import AnyCloudEvent
from cloudevents.http import CloudEvent
from jsonpath_ng import parse
from vanus.connect.source.custom import Message


class LabelMaker:
    def __init__(self, config=None, config_file=None):
        if config is None:
            if config_file is None:
                raise ValueError("config or config_file must be specified")
            with open(config_file) as f:
                config = json.load(f)

        patterns = self._build_patterns(config)
        matcher = Matcher.create_by_string(patterns, all_as_plain=True)
        if matcher is None:
            raise RuntimeError("Failed to create matcher")

        self._matcher = matcher

    def _build_patterns(self, config):
        buffer = StringIO()
        for label, detail in config.items():
            keys = detail.get("keys")
            if keys is None:
                continue
            for key in keys:
                buffer.write(f"{key}\t{label}\n")
        return buffer.getvalue()

    def label(self, text):
        labels = set()
        for match in self._matcher.finditer(text):
            labels.add(match[3])
        return list(labels)


class CloudEventLabelMaker:
    def __init__(self, source_path, target_path, *args, **kwargs):
        self._source_expr = parse(source_path)
        self._target_expr = parse(target_path)
        self._label_maker = LabelMaker(*args, **kwargs)

    def label(self, event: AnyCloudEvent) -> Optional[AnyCloudEvent]:
        data = event.get("data")
        if data is None:
            return None

        for match in self._source_expr.find(data):
            labels = self._label_maker.label(match.value)
            self._target_expr.update_or_create(data, labels)
            break

        return event


class HttpLabelMaker:
    def __init__(self, source_path, target_path, ce_source=None, ce_type=None, *args, **kwargs):
        if ce_source is None:
            ce_source = "https://github.com/vanus-labs/vanus-connect/connectors/source-labeling"
        if ce_type is None:
            ce_type = "ai.vanus.connect.source.labeling"

        self._source_expr = parse(source_path)
        self._target_expr = parse(target_path)
        self._label_maker = LabelMaker(*args, **kwargs)
        self._ce_source = ce_source
        self._ce_type = ce_type

    def label(self, msg: Message) -> CloudEvent:
        for match in self._source_expr.find(msg):
            labels = self._label_maker.label(match.value)
            self._target_expr.update_or_create(msg, labels)
            break

        attributes = {
            "source": self._ce_source,
            "type": self._ce_type,
        }
        event = CloudEvent.create(attributes, msg)

        return event
