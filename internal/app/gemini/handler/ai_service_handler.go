package handler

import (
	"github.com/haramurti/-KWU--food-web-be/internal/app/gemini/contract"
	"github.com/haramurti/-KWU--food-web-be/internal/app/gemini/dto"
	"github.com/haramurti/-KWU--food-web-be/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	aiService contract.AIService
}

func NewAIHandler(aiService contract.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

func (h *AIHandler) Chat(c *fiber.Ctx) error {
	var req dto.ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Message == "" {
		return response.BadRequest(c, "Message cannot be empty")
	}

	reply, err := h.aiService.Chat(req.History, req.Message)
	if err != nil {
		return response.InternalError(c, "Failed to get AI response: "+err.Error())
	}

	return response.OK(c, "Success", dto.ChatResponse{Reply: reply})
}
