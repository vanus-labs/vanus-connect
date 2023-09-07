# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from http import HTTPStatus

from cloudevents.http import from_http
from quart import Quart, ResponseReturnValue, request
from quart.views import View

from .knowledge_base import KnowledgeBase


class CloudEventsHandler(View):
    methods = ["POST"]
    init_every_request = True

    def __init__(self, base: KnowledgeBase) -> None:
        super().__init__()
        self._base = base

    async def dispatch_request(self, **kwargs) -> ResponseReturnValue:
        # Create a CloudEvent
        event = from_http(request.headers, await request.get_data())

        data = event.get_data()

        await self._base.load_document_from_texts(data)

        return "", HTTPStatus.NO_CONTENT


class KnowledgeBaseController:
    def __init__(self, kb: KnowledgeBase, app=None, **kwargs):
        self._kb = kb

        if app is None:
            app = Quart(__name__)

        self.app = app
        self._register_routes()

    def _register_routes(self):
        self.app.add_url_rule(
            "/",
            view_func=CloudEventsHandler.as_view("load-cloudevent", self._kb),
        )
