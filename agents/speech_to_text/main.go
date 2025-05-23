// To make this example work make sure you have speech.ogg file in the root of directory.
// You can use text to speech example to generate speech file.
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

	data, err := os.ReadFile("speech.mp3")
	if err != nil {
		panic(err)
	}

	result, err := factory.SpeechToText(openai.SpeechToTextParams{
		Model: goopenai.Whisper1,
	}).Execute(
		context.Background(),
		agenty.NewMessage(agenty.UserRole, agenty.VoiceKind, data),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(result.Content()))
}
