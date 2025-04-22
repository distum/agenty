package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/distum/agenty"
	"github.com/distum/agenty/providers/openai"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	result, err := factory.
		TextToText(openai.TextToTextParams{Model: goopenai.GPT4oMini}).
		SetPrompt("You are a helpful assistant that translates English to French").
		Execute(
			context.Background(),
			agenty.NewMessage(
				agenty.UserRole,
				agenty.TextKind,
				[]byte("I love programming."),
			),
		)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(result.Content()))
}
