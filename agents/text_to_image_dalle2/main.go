package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/distum/agenty"
	"github.com/distum/agenty/providers/openai"
)

func main() {
	provider := openai.New(openai.Params{
		Key: os.Getenv("OPENAI_API_KEY"),
	})

	result, err := provider.TextToImage(openai.TextToImageParams{
		Model:     "dall-e-2",
		ImageSize: "512x512",
		Quality:   "standard",
		Style:     "vivid",
	}).Execute(
		context.Background(),
		agenty.NewMessage(agenty.UserRole, agenty.TextKind, []byte("Halloween night at a haunted museum")),
	)
	if err != nil {
		panic(err)
	}

	if err := saveToDisk(result); err != nil {
		panic(err)
	}

	fmt.Println("Image has been saved!")
}

func saveToDisk(msg agenty.Message) error {
	r := bytes.NewReader(msg.Content())

	// for dall-e-3 use third party libraries due to lack of webp support in go stdlib
	imgData, format, err := image.Decode(r)
	if err != nil {
		return err
	}

	file, err := os.Create("example." + format)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		return err
	}

	return nil
}
