# Douyin Source

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
auth_code:
client_key:
client_secret:
EOF
```


| Name           | Required                           | Default | Description |
|:---------------|:-----------------------------------|:--------|:------------|
| auth_code      | YES                                |         | 抖音临时授权码     |
| client_key     | YES                                |         |             |
| client_secret  | YES                                |         |             |
