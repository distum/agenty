package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/distum/agenty"
	openAIProvider "github.com/distum/agenty/providers/openai"
	"github.com/sashabaranov/go-openai"
)

func main() { 
	imgBytes, err := os.ReadFile("example.png")
	if err != nil {
		panic(err)
	}

	result, err := openAIProvider.New(openAIProvider.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		ImageToText(openAIProvider.ImageToTextParams{Model: openai.GPT4o, MaxTokens: 300}).
		SetPrompt("describe what you see").
		Execute(
			context.Background(),
			agenty.NewMessage(agenty.UserRole, agenty.ImageKind, imgBytes),
		)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result.Content()))
}
