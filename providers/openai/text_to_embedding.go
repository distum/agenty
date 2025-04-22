package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/distum/agenty"
)

type EmbeddingModel = openai.EmbeddingModel

const AdaEmbeddingV2 EmbeddingModel = openai.AdaEmbeddingV2

type TextToEmbeddingParams struct {
	Model      EmbeddingModel
	Dimensions EmbeddingDimensions
}

type EmbeddingDimensions *int

func NewDimensions(v int) EmbeddingDimensions {
	return &v
}

func (p Provider) TextToEmbedding(params TextToEmbeddingParams) *agenty.Operation {
	var dimensions int
	if params.Dimensions != nil {
		dimensions = *params.Dimensions
	}

	return agenty.NewOperation(func(ctx context.Context, msg agenty.Message, cfg *agenty.OperationConfig) (agenty.Message, error) {
		// TODO: we have to convert string to model and then model to string. Can we optimize it?
		messages := append(cfg.Messages, msg)
		texts := make([]string, len(messages))

		for i, m := range messages {
			texts[i] = string(m.Content())
		}

		resp, err := p.client.CreateEmbeddings(
			ctx,
			openai.EmbeddingRequest{
				Input:      texts,
				Model:      params.Model,
				Dimensions: dimensions,
			},
		)
		if err != nil {
			return nil, err
		}

		vectors := make([]Embedding, len(resp.Data))
		for i, vector := range resp.Data {
			vectors[i] = vector.Embedding
		}

		bytes, err := EmbeddingToBytes(1536, vectors)
		if err != nil {
			return nil, fmt.Errorf("failed to convert embedding to bytes: %w", err)
		}

		// TODO: we have to convert []float32 to []byte. Can we optimize it?
		return agenty.NewMessage(agenty.AssistantRole, agenty.EmbeddingKind, bytes), nil
	})
}
