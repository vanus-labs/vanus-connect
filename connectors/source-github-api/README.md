---
title: GitHub API
---

# GitHub API Source

## Introduction

The GitHub API Source is a [Vanus Connector][vc] which aims to fetch contributors of target github organization.

The contributor's data is converted to:

```json
{
  "specversion": "1.0",
  "id": "026046e2-3cb0-4116-895e-c77877072dd2",
  "source": "https://github.com/apache",
  "datacontenttype": "application/json",
  "time": "2023-01-28T06:11:10.012579049Z",
  "data": {
  }
}
```

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
org_name: org_name
github_access_token: github_access_token
EOF
```

| Name                | Required | Default | Description                          |
|:--------------------|:---------|:--------|:-------------------------------------|
| target              | YES      |         | the target URL to send CloudEvents   |
| org_name            | YES      |         | github organization name, ex: apache |
| github_access_token | YES      |         | the github api access token          |

