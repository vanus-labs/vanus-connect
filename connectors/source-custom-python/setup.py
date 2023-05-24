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

from setuptools import find_namespace_packages, setup

if __name__ == "__main__":
    setup(
        name="vanus-connect-customsource",
        description="CustomSource of Vanus Connect.",
        author="Linkall Inc.",
        url="https://github.com/vanus-labs/vanus-connect",
        license="Apache License 2.0",
        packages=find_namespace_packages(include=["vanus.connect.customsource"]),
        classifiers=[
            "Operating System :: OS Independent",
            "Programming Language :: Python :: 3",
        ],
        install_requires=[
            "cloudevents",
            "httpx[http2]",
            "hypercorn",
            "quart",
        ],
        zip_safe=True,
    )
