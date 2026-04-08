package llm

type StreamEventType string

type TextDelta struct {
	Content string `json:"content"`
}

func (t TextDelta) String() string {
	return t.Content
}

const (
	TEXT_DELTA       StreamEventType = "text_delta"
	MESSAGE_COMPLETE StreamEventType = "message_complete"
	ERROR            StreamEventType = "error"
)

type TokenUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
	CachedTokens     int64 `json:"cached_tokens"`
}

func (t TokenUsage) Add(other TokenUsage) TokenUsage {
	return TokenUsage{
		PromptTokens:     t.PromptTokens + other.PromptTokens,
		CompletionTokens: t.CompletionTokens + other.CompletionTokens,
		TotalTokens:      t.TotalTokens + other.TotalTokens,
		CachedTokens:     t.CachedTokens + other.CachedTokens,
	}
}

type StreamEvent struct {
	Type         StreamEventType `json:"type"`
	TextDelta    *TextDelta      `json:"text_delta,omitempty"`
	Error        *string         `json:"error,omitempty"`
	FinishReason *string         `json:"finish_reason,omitempty"`
	Usage        *TokenUsage     `json:"usage,omitempty"`
}
