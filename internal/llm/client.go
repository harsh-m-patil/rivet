package llm

import (
	"context"
	"time"

	"github.com/openai/openai-go/v3"
)

type Messages []openai.ChatCompletionMessageParamUnion

type Result struct {
	Event StreamEvent
	Err   error
}

type Client struct {
	client     *openai.Client
	maxRetries int
}

func (c *Client) Get() *openai.Client {
	if c.client != nil {
		return c.client
	}

	client := openai.NewClient()
	c.client = &client
	c.maxRetries = 3
	return c.client
}

func (c *Client) Close() {
	c.client = nil
}

func (c *Client) ChatCompletion(
	ctx context.Context,
	messages Messages,
	stream bool,
) <-chan Result {
	ch := make(chan Result)
	ctx = context.Background()

	go func() {
		defer close(ch)

		for attempt := 0; attempt <= c.maxRetries; attempt++ {

			if stream {
				err := c.streamResponse(ctx, messages, ch)
				if err == nil {
					return
				}
				if !retryable(err) || attempt == c.maxRetries {
					errMsg := err.Error()
					ch <- Result{Event: StreamEvent{
						Type:  ERROR,
						Error: &errMsg,
					}}
					return
				}
			} else {
				ev, err := c.nonStreamResponse(ctx, messages)
				if err == nil {
					ch <- Result{Event: ev}
					return
				}
				if !retryable(err) || attempt == c.maxRetries {
					errMsg := err.Error()
					ch <- Result{Event: StreamEvent{
						Type:  ERROR,
						Error: &errMsg,
					}}
					return
				}
			}

			time.Sleep(time.Duration(1<<attempt) * time.Second)
		}
	}()

	return ch
}

func retryable(err error) bool {
	// TODO: implement this later
	return true
}

func (c *Client) streamResponse(
	ctx context.Context,
	messages Messages,
	ch chan<- Result,
) error {
	stream := c.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    openai.ChatModelGPT4o,
	})

	defer stream.Close()

	var finishReason string
	var usage *TokenUsage

	for stream.Next() {
		chunk := stream.Current()

		if chunk.Usage.TotalTokens > 0 {
			usage = &TokenUsage{
				PromptTokens:     chunk.Usage.PromptTokens,
				CompletionTokens: chunk.Usage.CompletionTokens,
				TotalTokens:      chunk.Usage.TotalTokens,
			}
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		choice := chunk.Choices[0]

		if choice.FinishReason != "" {
			finishReason = choice.FinishReason
		}

		if choice.Delta.Content != "" {
			ch <- Result{
				Event: StreamEvent{
					Type: TEXT_DELTA,
					TextDelta: &TextDelta{
						Content: choice.Delta.Content,
					},
				},
			}
		}
	}

	if err := stream.Err(); err != nil {
		return err
	}

	ch <- Result{
		Event: StreamEvent{
			Type:         MESSAGE_COMPLETE,
			FinishReason: &finishReason,
			Usage:        usage,
		},
	}

	return nil
}

func (c *Client) nonStreamResponse(
	ctx context.Context,
	messages Messages,
) (StreamEvent, error) {

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModelGPT5_4Nano,
		Messages: messages,
	})
	if err != nil {
		return StreamEvent{}, err
	}

	choice := resp.Choices[0]

	var textDelta *TextDelta
	if choice.Message.Content != "" {
		textDelta = &TextDelta{Content: choice.Message.Content}
	}

	var usage *TokenUsage
	if resp.Usage.TotalTokens > 0 {
		usage = &TokenUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}
	}

	return StreamEvent{
		Type:         MESSAGE_COMPLETE,
		TextDelta:    textDelta,
		FinishReason: &choice.FinishReason,
		Usage:        usage,
	}, nil
}
