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

- config when list_type is "org"
```shell
cat << EOF > config.yml
target: http://localhost:31081
github_access_token: github_access_token
api_type: contributor
list_type: org
organizations:
  - apache
  - google
EOF
```

- config when list_type is "user"
```shell
cat << EOF > config.yml
target: http://localhost:31081
github_access_token: github_access_token
api_type: contributor
list_type: user
user_list:
  - u1
  - u2
EOF
```

- config when api_type is "pr"
```shell
cat << EOF > config.yml
target: http://localhost:31081
github_access_token: github_access_token
api_type: pr
pr_configs:
  - organization: apache
    repo: spark
    user_list: 
      - user1
      - user2
  - organization: microsoft
    repo: vcpkg
    user_list: 
      - user3
      - user4
EOF
```

Config
---
| Name                | Required                           | Default | Description                        |
|:--------------------|:-----------------------------------|:--------|:-----------------------------------|
| target              | YES                                |         | the target URL to send CloudEvents |
| github_access_token | YES                                |         | the github api access token        |
| api_type            | YES                                |         | "contributor" or "pr"              |
| list_type           | YES when api_type is "contributor" |         | "org" or "user"                    |
| organizations       | YES when list_type is "org"        |         | organization arrays                |
| user_list           | YES when list_type is "user"       |         | uid arrays                         |
| pr_configs          | YES when api_type is "pr"          |         | PRConfig arrays                    |  

PRConfig
---
| Name                    |  Description                        |
|:------------------------|:------------------------------------|
| organization            |  organization                       |
| repo                    |  repository                         |
| user_list               |  array of contributor id            | 
