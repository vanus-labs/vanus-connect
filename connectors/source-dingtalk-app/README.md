# dingtalk-app source

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
vanus_ai_app_id:
vanus_ai_url:
dingtalk_app_key: 
dingtalk_app_secret: 
EOF
```

Config
---
| Name                | Required | Type   | Default | Description      |
|:--------------------|:---------|:-------|:--------|:-----------------|
| vanus_ai_app_id     | YES      | String |         | VanusAI 应用ID     |
| vanus_ai_url        | YES      | String |         | VanusAI URL      |
| dingtalk_app_key    | YES      | String |         | 钉钉应用 AgentKey    |
| dingtalk_app_secret | YES      | String |         | 钉钉应用 AgentSecret |