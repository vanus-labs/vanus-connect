# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import argparse

import aiorun
from vanus.connect.cdk import build_pipeline
from vanus.connect.cloudevents import CloudEventSource

from .crawler import CrawlerService


def main():
    parser = argparse.ArgumentParser(
        prog="webpage-service", description="vanus connect service webpage", epilog="Linkall Inc."
    )
    parser.add_argument("--port", default=8080, type=int, help="the webpage service port")
    args = parser.parse_args()

    aiorun.run(
        build_pipeline(CloudEventSource(args.port))
        .then(
            CrawlerService(),
        )
        .start()
    )


if __name__ == "__main__":
    main()
