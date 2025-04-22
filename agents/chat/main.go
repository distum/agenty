package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/distum/agenty"
	"github.com/distum/agenty/providers/openai"
)

func main() {
    OPENAI_API_KEY := "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"
    // fmt.Print(OPENAI_API_KEY + "\n")
	assistant := openai.
		New(openai.Params{Key: os.Getenv(OPENAI_API_KEY)}).
		TextToText(openai.TextToTextParams{Model: "nvidia/llama-3.1-nemotron-ultra-253b-v1:free"}).
		SetPrompt("You are helpful assistant.")

	messages := []agenty.Message{}
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Print("User: ")

		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input := agenty.NewTextMessage(agenty.UserRole, text)
		answer, err := assistant.SetMessages(messages).Execute(ctx, input)
		if err != nil {
			panic(err)
		}

		fmt.Println("Assistant:", string(answer.Content()))

		messages = append(messages, input, answer)

        fmt.Println("--------------------------------------------------------")
	}
}
