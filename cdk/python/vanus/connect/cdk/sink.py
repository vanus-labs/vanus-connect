# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from abc import ABC, abstractmethod
from typing import Optional, TypeVar

from cloudevents.abstract import CloudEvent

AnySink = TypeVar("AnySink", bound="Sink")


class Sink(ABC):
    def __init__(self):
        super().__init__()

    async def start(self):
        """start the sink"""

    @abstractmethod
    async def on_event(self, event: CloudEvent) -> Optional[CloudEvent]:
        """receive event from source"""
