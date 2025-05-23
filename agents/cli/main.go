package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/distum/agenty"
	"github.com/distum/agenty/providers/openai"
)

// usage example: go to the repo root and execute
// go run examples/cli/main.go -prompt "You are professional translator, translate everything you see to Russian" -model "gpt-4o-mini" -maxTokens=1000 "I love winter"
func main() {
	provider := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	temp := flag.Float64("temp", 0.0, "Temperature value")
	maxTokens := flag.Int("maxTokens", 0, "Maximum number of tokens")
	model := flag.String("model", "gpt-4o-mini", "Model name")
	prompt := flag.String("prompt", "You are a helpful assistant", "System message")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("content argument is required")
		os.Exit(1)
	}
	content := args[0]

	result, err := provider.
		TextToText(openai.TextToTextParams{
			Model:       *model,
			Temperature: openai.Temperature(float32(*temp)),
			MaxTokens:   *maxTokens,
		}).
		SetPrompt(*prompt).
		Execute(context.Background(), agenty.NewMessage(agenty.UserRole, agenty.TextKind, []byte(content)))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(result.Content())
}
