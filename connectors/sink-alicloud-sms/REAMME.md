# Aliyun SMS Sink

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
access_key_id:
access_key_secret:
sign_name: 
template_code: 
template_param: 
EOF
```


| Name              | Required | Default | Description             |
|:------------------|:---------|:--------|:------------------------|
| access_key_id     | YES      |         | aliyun AccessKey ID     |
| access_key_secret | YES      |         | aliyun AccessKey Secret |
| sign_name         | YES      |         | 短信签名名称                  |
| template_code     | YES      |         | 短信模板CODE                |
| template_param    | Optional |         | 短信模板变量对应的实际值            |

