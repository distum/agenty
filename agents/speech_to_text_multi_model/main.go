// To make this example work make sure you have speech.ogg file in the root of directory
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

type Saver []agenty.Message

func (s *Saver) Save(input, output agenty.Message, _ *agenty.OperationConfig) {
	*s = append(*s, output)
}

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	// step 1
	hear := factory.
		SpeechToText(openai.SpeechToTextParams{
			Model: goopenai.Whisper1,
		})

	// step2
	translate := factory.
		TextToText(openai.TextToTextParams{
			Model:       "gpt-4o-mini",
			Temperature: openai.Temperature(0.5),
		}).
		SetPrompt("translate to russian")

	// step 3
	uppercase := factory.
		TextToText(openai.TextToTextParams{
			Model:       "gpt-4o-mini",
			Temperature: openai.Temperature(1),
		}).
		SetPrompt("uppercase every letter of the text")

	saver := Saver{}

	sound, err := os.ReadFile("speech.mp3")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	speechMsg := agenty.NewMessage(agenty.UserRole, agenty.VoiceKind, sound)

	_, err = agenty.NewProcess(
		hear,
		translate,
		uppercase,
	).Execute(ctx, speechMsg, saver.Save)
	if err != nil {
		panic(err)
	}

	for _, msg := range saver {
		fmt.Println(string(msg.Content()))
	}
}
