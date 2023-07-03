package vanus_ai

type Config struct {
	URL   string `json:"url" yaml:"url" validate:"required"`
	AppID string `json:"app_id" yaml:"app_id" validate:"required"`
}
