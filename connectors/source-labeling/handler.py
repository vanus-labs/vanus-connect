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


from vanus.connect.source.custom import EventHandler
from vanus.connect.source.labeling import CloudEventLabelMaker


def create_handler(*args, **kwargs) -> EventHandler:
    label_maker = CloudEventLabelMaker(*args, **kwargs)
    return label_maker.label
