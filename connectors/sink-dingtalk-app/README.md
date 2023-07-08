# dingtalk-app sink

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
dingtalk_app_key: 
dingtalk_app_secret: 
EOF
```

Config
---
| Name                  | Required | Type   | Default | Description      |
|:----------------------|:---------|:-------|:--------|:-----------------|
| dingtalk_app_key      | YES      | String |         | 钉钉应用 AgentKey    |
| dingtalk_app_secret   | YES      | String |         | 钉钉应用 AgentSecret |