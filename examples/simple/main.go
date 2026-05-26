package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zraisan/gotique"
	"github.com/zraisan/gotique/models"
	"github.com/zraisan/gotique/providers/openai"
)

func main() {
	provider := openai.New(os.Getenv("OPENAI_API_KEY"))
	model := models.New("gpt-4.1-mini", provider)

	agent := gotique.NewAgent(
		gotique.WithName("assistant"),
		gotique.WithModel(model),
	)

	response, err := agent.Run(context.Background(), "Say hello in one sentence.")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
