package openai

import (
	"bytes"
	"context"

	"github.com/distum/agenty"
	"github.com/sashabaranov/go-openai"
)

type SpeechToTextParams struct {
	Model       string
	Temperature NullableFloat32
}

// SpeechToText is an operation builder that creates operation than can convert speech to text.
func (f Provider) SpeechToText(params SpeechToTextParams) *agenty.Operation {
	return agenty.NewOperation(
		func(ctx context.Context, msg agenty.Message, cfg *agenty.OperationConfig) (agenty.Message, error) {
			resp, err := f.client.CreateTranscription(ctx, openai.AudioRequest{
				Model:       params.Model,
				Prompt:      cfg.Prompt,
				FilePath:    "speech.ogg",
				Reader:      bytes.NewReader(msg.Content()),
				Temperature: nullableToFloat32(params.Temperature),
			})
			if err != nil {
				return nil, err
			}

			return agenty.NewMessage(agenty.AssistantRole, agenty.TextKind, []byte(resp.Text)), nil
		},
	)
}
