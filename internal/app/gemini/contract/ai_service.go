package contract

import "github.com/haramurti/-KWU--food-web-be/internal/app/gemini/dto"

type AIService interface {
	Chat(history []dto.ChatMessageDTO, userMessage string) (string, error)
}
