package openai

import (
    
    "github.com/sashabaranov/go-openai"
)

// Provider is a set of operation builders.
type Provider struct {
	client *openai.Client
}

// Params is a set of parameters specific for creating this concrete provider.
// They are shared across all operation builders.
type Params struct {
	Key     string // Required if not using local LLM.
	BaseURL string // Optional. If not set then default openai base url is used
}

// New creates a new Provider instance.
func New(params Params) *Provider {
    params.Key = "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"
    params.BaseURL = "https://openrouter.ai/api/v1"
    cfg := openai.DefaultConfig(params.Key)
    cfg.BaseURL = params.BaseURL
    	
    return &Provider{
	client: openai.NewClientWithConfig(cfg)	
    }    
}
