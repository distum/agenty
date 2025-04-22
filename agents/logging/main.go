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
	params := openai.TextToTextParams{Model: "gpt-4o-mini"}

	_, err := agenty.NewProcess(
		factory.TextToText(params).SetPrompt("explain what that means"),
		factory.TextToText(params).SetPrompt("translate to russian"),
		factory.TextToText(params).SetPrompt("replace all spaces with '_'"),
	).
		Execute(
			context.Background(),
			agenty.NewMessage(agenty.UserRole, agenty.TextKind, []byte("Kazakhstan alga!")),
			Logger,
		)

	if err != nil {
		panic(err)
	}
}

func Logger(input, output agenty.Message, cfg *agenty.OperationConfig) {
	fmt.Printf(
		"in: %v\nprompt: %v\nout: %v\n\n",
		string(input.Content()),
		cfg.Prompt,
		string(output.Content()),
	)
}
