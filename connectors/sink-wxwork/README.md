# wxwork sink

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
wework_corp_id: 
wework_agent_id: 
wework_agent_secret:
EOF
```

Config
---
| Name                    | Required | Default                              | Description           |
|:------------------------|:---------|:-------------------------------------|:----------------------|
| wework_corp_id          | YES      |                                      | 企业微信 企业ID             |
| wework_agent_id         | YES      |                                      | 企业微信 AgentId          |
| wework_agent_secret     | YES      |                                      | 企业微信 AgentSecret      |