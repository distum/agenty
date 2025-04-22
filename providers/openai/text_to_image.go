package openai

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/distum/agenty"
	"github.com/sashabaranov/go-openai"
)

type TextToImageParams struct {
	Model     string
	ImageSize string
	Quality   string
	Style     string
}

// TextToImage is an operation builder that creates operation than can convert text to image.
func (p Provider) TextToImage(params TextToImageParams) *agenty.Operation {
	return agenty.NewOperation(
		func(ctx context.Context, msg agenty.Message, cfg *agenty.OperationConfig) (agenty.Message, error) {
			reqBase64 := openai.ImageRequest{
				Prompt:         fmt.Sprintf("%s\n\n%s", cfg.Prompt, string(msg.Content())),
				Size:           params.ImageSize,
				ResponseFormat: openai.CreateImageResponseFormatB64JSON,
				N:              1, // DALL·E-3 only support n=1, for other models support needed
				Model:          params.Model,
				Quality:        params.Quality,
				Style:          params.Style,
			}

			respBase64, err := p.client.CreateImage(ctx, reqBase64)
			if err != nil {
				return nil, err
			}

			imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
			if err != nil {
				return nil, err
			}

			return agenty.NewMessage(agenty.AssistantRole, agenty.ImageKind, imgBytes), nil
		},
	)
}
