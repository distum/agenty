package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/distum/agenty"
	"github.com/distum/agenty/providers/openai"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	resultMsg, err := factory.
		TextToText(openai.TextToTextParams{Model: "gpt-4o-mini"}).
		SetPrompt(
			"You are a helpful assistant that translates %s to %s",
			"English", "French",
		).
		Execute(
			context.Background(),
			agenty.NewMessage(agenty.UserRole, agenty.TextKind, []byte("I love programming.")),
		)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(resultMsg.Content()))
}
