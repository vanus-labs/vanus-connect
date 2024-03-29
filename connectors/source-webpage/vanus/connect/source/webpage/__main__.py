# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

import argparse

import aiorun
import yaml
from vanus.connect.cdk import build_pipeline
from vanus.connect.cloudevents import CloudEventSink
from vanus.connect.service.webpage import CrawlerService

from ..timer import TimerSource


def main():
    parser = argparse.ArgumentParser(
        prog="webpage-source", description="vanus connect source webpage", epilog="Linkall Inc."
    )
    parser.add_argument("--config", help="the webpage source config")
    args = parser.parse_args()

    with open(args.config) as f:
        config = yaml.safe_load(f)

    aiorun.run(
        build_pipeline(
            TimerSource(float(config["interval"]), {"url": config["url"]}),
        )
        .call(
            CrawlerService(),
        )
        .then(
            CloudEventSink(config["target"]),
        )
        .start()
    )


if __name__ == "__main__":
    main()
