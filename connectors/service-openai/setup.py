# SPDX-FileCopyrightText: 2023 Linkall Inc.
#
# SPDX-License-Identifier: Apache-2.0

from setuptools import find_namespace_packages, setup

if __name__ == "__main__":
    setup(
        name="vanus-connect-openaiservice",
        description="OpenAI Service of Vanus Connect.",
        author="Linkall Inc.",
        url="https://github.com/vanus-labs/vanus-connect",
        license="Apache License 2.0",
        packages=find_namespace_packages(
            include=[
                "vanus.connect.service.openai",
            ]
        ),
        classifiers=[
            "Operating System :: OS Independent",
            "Programming Language :: Python :: 3",
        ],
        install_requires=[
            "openai",
            "vanus-cdk",
        ],
        # extras_require={
        #     "run": [
        #         "aiorun",
        #     ],
        # },
        zip_safe=True,
    )
