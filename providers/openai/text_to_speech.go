package openai

import (
	"context"
	"io"

	"github.com/distum/agenty"
	"github.com/sashabaranov/go-openai"
)

type TextToSpeechParams struct {
	Model          string
	ResponseFormat string
	Speed          float64
	Voice          string
}

// TextToSpeech is an operation builder that creates operation than can convert text to speech.
func (f Provider) TextToSpeech(params TextToSpeechParams) *agenty.Operation {
	return agenty.NewOperation(
		func(ctx context.Context, msg agenty.Message, cfg *agenty.OperationConfig) (agenty.Message, error) {
			resp, err := f.client.CreateSpeech(ctx, openai.CreateSpeechRequest{
				Model:          openai.SpeechModel(params.Model),
				Input:          string(msg.Content()),
				Voice:          openai.SpeechVoice(params.Voice),
				ResponseFormat: openai.SpeechResponseFormat(params.ResponseFormat),
				Speed:          params.Speed,
			})
			if err != nil {
				return nil, err
			}

			bb, err := io.ReadAll(resp)
			if err != nil {
				return nil, err
			}

			return agenty.NewMessage(agenty.AssistantRole, agenty.VoiceKind, bb), nil
		},
	)
}
