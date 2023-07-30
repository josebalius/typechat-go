package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/josebalius/typechat-go"
	openaiadapter "github.com/josebalius/typechat-go/adapters/openai"
	"github.com/sashabaranov/go-openai"
)

type Sentiment int

const (
	Positive Sentiment = iota
	Negative
	Neutral
)

func (s Sentiment) String() string {
	switch s {
	case Positive:
		return "positive"
	case Negative:
		return "negative"
	case Neutral:
		return "neutral"
	}

	return ""
}

type SentimentAnalysis struct {
	Sentiment  Sentiment
	Confidence float64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ./examples/openai-sentiment '<prompt>'")
		return
	}
	message := strings.Join(os.Args[1:], " ")

	ctx := context.Background()
	token := os.Getenv("OPENAI_TOKEN")
	client := openai.NewClient(token)
	model := openaiadapter.NewClient(client, openai.GPT3Dot5Turbo)

	prompt := typechat.NewPrompt[SentimentAnalysis](model, message)
	analysis, err := prompt.Execute(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sentiment: %s\n", analysis.Sentiment)
}
