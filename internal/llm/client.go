package llm

import (
	"context"

	"github.com/openai/openai-go/v3"
)

func GetAnswer(prompt string) {
	client := openai.NewClient()
	ctx := context.Background()

	stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are Rivet. A terminal based minimal coding agent"),
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT5_2,
	})

	for stream.Next() {
		evt := stream.Current()
		if len(evt.Choices) > 0 {
			print(evt.Choices[0].Delta.Content)
		}
	}
	println()

	if err := stream.Err(); err != nil {
		panic(err.Error())
	}

}
