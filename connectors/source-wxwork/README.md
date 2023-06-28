# wxwork source

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
vanus_ai_app_id: 
vanus_ai_url: 
wework_corp_id: 
wework_agent_id: 
wework_agent_secret:
wework_token:
wework_encoding_aes_key:
EOF
```

Config
---
| Name                    | Required | Default                              | Description        |
|:------------------------|:---------|:-------------------------------------|:-------------------|
| vanus_ai_app_id         | YES      |                                      | vanus-ai 应用ID      |
| vanus_ai_url            | YES      | 各环境不同，例如海外线上 https://app.ai.vanus.ai | vanus-ai 地址        |
| wework_corp_id          | YES      |                                      | 企业微信 企业ID          |
| wework_agent_id         | YES      |                                      | 企业微信 AgentId       |
| wework_agent_secret     | YES      |                                      | 企业微信 AgentSecret   |
| wework_token            | YES      |                                      | 企业微信 Token         |
| wework_encoding_aes_key | YES      |                                      | 企业微信 企业微信 AgentId  |
