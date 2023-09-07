# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import asyncio
import functools
from typing import Any, Awaitable, Callable


def ensure_async(func: Callable[..., Any], executor=None) -> Callable[..., Awaitable[Any]]:
    if asyncio.iscoroutinefunction(func):
        return func
    else:

        async def _wrapper(*args, **kwargs):
            return await asyncio.get_running_loop().run_in_executor(executor, functools.partial(func, *args, **kwargs))

        return _wrapper
