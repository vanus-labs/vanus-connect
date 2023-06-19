# Aliyun SMS Sink

## Quick Start

### Create the config file

```shell
cat << EOF > config.yml
access_key_id:
access_key_secret:
sign_name: 
phone_numbers: $.data.phonenumbers
template_code: 
template_param: 
  - key: code1
    value: $.data.code1
  - key: code2
    value: 333333
EOF
```

Config
---
| Name                 | Required | Default | Description                                                 |
|:---------------------|:---------|:--------|:------------------------------------------------------------|
| access_key_id        | YES      |         | aliyun AccessKey ID                                         |
| access_key_secret    | YES      |         | aliyun AccessKey Secret                                     |
| sign_name            | YES      |         | 短信签名名称                                                      |
| phone_numbers        | YES      |         | 手机号，支持常量(131xxx,186xxx)，或指定field从event中获取(例如：$.data.phones) |
| template_code        | YES      |         | 短信模板CODE                                                    |
| template_param       | YES      |         | TemplateKV数组                                                |

TemplateKV
---
| Name  | Required | Default | Description                                         |
|:------|:---------|:--------|:----------------------------------------------------|
| key   | YES      |         | 短信模板变量                                              |
| value | YES      |         | 短信模板变量对应的实际值，支持常量，或指定field从event中获取(例如：$.data.code) |
