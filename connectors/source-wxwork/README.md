# wxwork source

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
vanus_ai_app_id: 
wework_corp_id: 
wework_agent_id: 
wework_agent_secret:
wework_token:
wework_encoding_aes_key:
EOF
```

Config
---
| Name                    | Required | Type   | Description            |
|:------------------------|:---------|:-------|:-----------------------|
| vanus_ai_app_id         | YES      | String | vanus-ai 应用ID          |
| wework_corp_id          | YES      | String | 企业微信 企业ID              |
| wework_agent_id         | YES      | Int    | 企业微信 AgentId           |
| wework_agent_secret     | YES      | String | 企业微信 AgentSecret       |
| wework_token            | YES      | String | 企业微信 Token             |
| wework_encoding_aes_key | YES      | String | 企业微信 EncodingAESKey    |
