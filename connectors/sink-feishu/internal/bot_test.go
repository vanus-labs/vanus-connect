package internal

import (
	"github.com/go-resty/resty/v2"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	m := `{
  "config": {
    "wide_screen_mode": true
  },
  "header": {
    "template": "<$.data.color>",
    "title": {
      "content": "Minimax告警消息",
      "tag": "plain_text"
    }
  },
  "elements": [
    {
      "tag": "div",
      "text": {
        "content": "名称: <$.data.body.commonLabels.alertname>\n级别: <$.data.body.commonAnnotations.level>\n状态: <$.data.body.status>\n集群: <$.data.body.commonLabels.cluster>",
        "tag": "plain_text"
      }
    },
    {
      "tag": "div",
      "text": {
        "content": "以下共有 <$.data.alerts_count> 条告警",
        "tag": "plain_text"
      }
    },
    {
      "tag": "div",
      "text": {
        "content": "<$.data.alerts_text>",
        "tag": "plain_text"
      }
    },
    {
      "tag": "hr"
    },
    {
      "elements": [
        {
          "content": "来自 <source>",
          "tag": "plain_text"
        }
      ],
      "tag": "note"
    }
  ]
}`
	tt := time.Now()
	payload := map[string]interface{}{
		"timestamp": tt.Unix(),
		"msg_type":  interactiveMessage,
		"card":      m,
	}
	r, e := resty.New().R().SetBody(payload).Post("")
	println(r)
	println(e)
}
