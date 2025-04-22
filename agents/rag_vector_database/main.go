package main

import (
	"context"
	"encoding/json"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/distum/agenty"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"

	"github.com/distum/agenty/providers/openai"
)

// natural langauge query -> weaviate RAG -> speech
func main() {
	openAPIKey := os.Getenv("OPENAI_API_KEY")

	ctx := context.Background()

	client, err := prepareDB(openAPIKey, ctx)
	if err != nil {
		panic(err)
	}

	factory := openai.New(openai.Params{Key: openAPIKey})
	retrieve := RAGoperation(client)
	summarize := factory.TextToText(openai.TextToTextParams{Model: "gpt-4o-mini"}).SetPrompt("summarize")
	voice := factory.TextToSpeech(openai.TextToSpeechParams{
		Model: "tts-1", ResponseFormat: "mp3", Speed: 1, Voice: "onyx",
	})

	result, err := agenty.NewProcess(
		retrieve,
		summarize,
		voice,
	).Execute(ctx, agenty.NewMessage(agenty.UserRole, agenty.TextKind, []byte("programming")))
	if err != nil {
		panic(err)
	}

	if err := saveToDisk(result); err != nil {
		panic(err)
	}
}

// RAGoperation retrieves relevant objects from vector store and builds a text message to pass further to the process
func RAGoperation(client *weaviate.Client) *agenty.Operation {
	return agenty.NewOperation(func(ctx context.Context, msg agenty.Message, po *agenty.OperationConfig) (agenty.Message, error) {
		input := string(msg.Content())

		result, err := client.GraphQL().Get().
			WithClassName("Records").
			WithFields(graphql.Field{Name: "content"}).
			WithNearText(
				client.GraphQL().
					NearTextArgBuilder().
					WithConcepts(
						[]string{input},
					),
			).
			WithLimit(10).
			Do(ctx)
		if err != nil {
			panic(err)
		}

		var content string
		for _, obj := range result.Data {
			bb, err := json.Marshal(&obj)
			if err != nil {
				return nil, err
			}
			content += string(bb)
		}

		return agenty.NewMessage(
			agenty.AssistantRole,
			agenty.TextKind,
			[]byte(content),
		), nil
	})
}

func prepareDB(openAPIKey string, ctx context.Context) (*weaviate.Client, error) {
	client, err := weaviate.NewClient(weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
		Headers: map[string]string{
			"X-OpenAI-Api-Key": openAPIKey,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := client.Schema().AllDeleter().Do(ctx); err != nil {
		return nil, err
	}

	classObj := &models.Class{
		Class:      "Records",
		Vectorizer: "text2vec-openai",
		ModuleConfig: map[string]interface{}{
			"text2vec-openai":   map[string]interface{}{},
			"generative-openai": map[string]interface{}{},
		},
	}
	if err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background()); err != nil {
		return nil, err
	}

	if _, err := client.Batch().ObjectsBatcher().WithObjects(data...).Do(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func saveToDisk(msg agenty.Message) error {
	file, err := os.Create("speech.mp3")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(msg.Content())
	if err != nil {
		return err
	}

	return nil
}
