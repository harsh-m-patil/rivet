package llm

import (
	"context"

	"github.com/openai/openai-go/v3"
)

func GetAnswer(prompt string) {
	client := openai.NewClient()

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are Rivet. A terminal based minimal coding agent"),
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT5_2,
	})
	if err != nil {
		panic(err)
	}

	println(chatCompletion.Choices[0].Message.Content)
}
