package dto

type ChatMessageDTO struct {
	Role    string `json:"role"` // "user" | "assistant"
	Content string `json:"content"`
}

type ChatRequest struct {
	Message string           `json:"message"`
	History []ChatMessageDTO `json:"history"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}
