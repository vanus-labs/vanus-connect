# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from setuptools import find_namespace_packages, setup

if __name__ == "__main__":
    setup(
        name="vanus-connect-webpageservice",
        description="Webpage Service of Vanus Connect.",
        author="Linkall Inc.",
        url="https://github.com/vanus-labs/vanus-connect",
        license="Apache License 2.0",
        packages=find_namespace_packages(
            include=[
                "vanus.connect.service.webpage",
            ]
        ),
        classifiers=[
            "Operating System :: OS Independent",
            "Programming Language :: Python :: 3",
        ],
        install_requires=[
            "aiorun",
            "beautifulsoup4",
            "bs2json",
            "httpx[http2]",
            "lxml",
            "vanus-cdk",
        ],
        zip_safe=True,
    )
