#!/usr/bin/env python3
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

import argparse
import importlib
import importlib.machinery
import importlib.util
import json
import os.path
import sys
from typing import Any, Dict, List

from .run import run_source
from .source import AsyncEventHandler, SyncEventHandler


def _resolve_handle(symbol: str, base=None):
    if ":" in symbol:
        mod, _, attr = symbol.rpartition(":")
    else:
        mod, attr = None, symbol
    if mod is not None:
        if base is None:
            module = importlib.import_module(mod)
        else:
            spec = importlib.machinery.PathFinder.find_spec(
                ".".join([base.__package__, mod]), base.__spec__.submodule_search_locations
            )
            if spec is None:
                raise RuntimeError("Failed to load spec")
            module = importlib.util.module_from_spec(spec)
            if spec.loader is None:
                raise RuntimeError("No loader for spec")
            spec.loader.exec_module(module)
    else:
        module = base
    return getattr(module, attr)


def resolve_handle(handle_spec: str):
    if "#" in handle_spec:
        path, _, symbol = handle_spec.rpartition("#")
    else:
        return _resolve_handle(handle_spec)

    if path == "":
        path = "."

    if os.path.isdir(path):
        path = os.path.join(path, "__init__.py")

    loader = importlib.machinery.SourceFileLoader("custom_handle", os.path.abspath(path))
    spec = importlib.util.spec_from_loader("custom_handle", loader)
    if spec is None:
        raise RuntimeError("Failed to load spec")

    module = importlib.util.module_from_spec(spec)
    if module is None:
        raise RuntimeError("Failed to load module")

    if spec.loader is None:
        raise RuntimeError("No loader for spec")

    sys.modules["custom_handle"] = module
    spec.loader.exec_module(module)

    return _resolve_handle(symbol, module)


def create_handler(handle_spec, *args, **kwargs) -> SyncEventHandler:
    factory = resolve_handle(handle_spec)
    handle = factory(*args, **kwargs)
    return handle


def create_ahandler(handle_spec, *args, **kwargs) -> AsyncEventHandler:
    factory = resolve_handle(handle_spec)
    handle = factory(*args, **kwargs)
    return handle


def main():
    parser = argparse.ArgumentParser(
        prog="custom-source", description="vanus connect customsource", epilog="Linkall Inc."
    )
    parser.add_argument("handle", help="The handle to process events as [path.to.package#]path.to.module:function.path")
    parser.add_argument("--name", help="the source name", required=False)
    parser.add_argument("--port", help="the source port", default=3000, type=int, required=False)
    parser.add_argument("--sink-endpoint", help="the sink endpoint")
    parser.add_argument("--handle-async", help="the flag to indicate handle is async", required=False)
    parser.add_argument("--handle-args", help="the handle factory args", required=False)
    parser.add_argument("--handle-kwargs", help="the handle factory args", required=False)
    args = parser.parse_args()

    handle_args: List[Any] = list()
    if args.handle_args is not None:
        handle_args = json.loads(args.handle_args)

    handle_kwargs: Dict[str, Any] = dict()
    if args.handle_kwargs is not None:
        handle_kwargs = json.loads(args.handle_kwargs)

    kwargs = {"name": args.name}

    if args.handle_async:
        kwargs["async_handler"] = create_ahandler(args.handle, *handle_args, **handle_kwargs)
    else:
        kwargs["sync_handler"] = create_handler(args.handle, *handle_args, **handle_kwargs)

    run_source(args.port, args.sink_endpoint, **kwargs)


if __name__ == "__main__":
    main()
