package vanus_ai

import "fmt"

type Config struct {
	URL   string `json:"url" yaml:"url" `
	AppID string `json:"app_id" yaml:"app_id"`
}

func (c Config) Validate() error {
	if c.URL == "" || c.AppID == "" {
		return fmt.Errorf("vanus-ai url or appid is empty")
	}
	return nil
}
