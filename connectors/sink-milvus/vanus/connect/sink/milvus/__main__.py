#!/usr/bin/env python3
# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import argparse

import yaml
from hypercorn.config import Config
from hypercorn.middleware import DispatcherMiddleware

from .knowledge_base import KnowledgeBase
from .server import KnowledgeBaseController


def _run(app, config):
    import asyncio

    from hypercorn.asyncio import serve

    try:
        import uvloop

        uvloop.install()
    except ImportError:
        pass

    dispatcher_app = DispatcherMiddleware(
        {
            "/": app,
        }
    )
    asyncio.run(serve(dispatcher_app, config))


def run_app(port, app):
    config = Config()
    config.bind = [f"0.0.0.0:{port}"]
    _run(app, config)


def run_knowledge_base_service(port=8080, **kwargs):
    kb = KnowledgeBase(**kwargs)

    srv = KnowledgeBaseController(kb, **kwargs)
    run_app(port, srv.app)


def main():
    parser = argparse.ArgumentParser(
        prog="vanus-sink-milvus", description="Vanus connect sink milvus Service", epilog="Linkall Inc."
    )
    parser.add_argument("--config", help="the path of configuration file")
    args = parser.parse_args()

    if args.config:
        with open(args.config) as f:
            config = yaml.safe_load(f)
    else:
        config = dict()

    run_knowledge_base_service(**config)


if __name__ == "__main__":
    main()
