package handler

import (
	"github.com/haramurti/-KWU--food-web-be/internal/app/payment/contract"
	"github.com/haramurti/-KWU--food-web-be/internal/app/payment/dto"
	"github.com/haramurti/-KWU--food-web-be/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	paymentService contract.PaymentService
	webhookToken   string
}

func NewPaymentHandler(paymentService contract.PaymentService, webhookToken string) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		webhookToken:   webhookToken,
	}
}

func (h *PaymentHandler) CreateInvoice(c *fiber.Ctx) error {
	var req dto.CreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.ExternalID == "" {
		return response.BadRequest(c, "external_id is required")
	}
	if req.Amount <= 0 {
		return response.BadRequest(c, "amount must be greater than 0")
	}
	if req.PayerEmail == "" {
		return response.BadRequest(c, "payer_email is required")
	}

	result, err := h.paymentService.CreateInvoice(req)
	if err != nil {
		return response.InternalError(c, "Failed to create invoice: "+err.Error())
	}

	return response.OK(c, "Invoice created", result)
}

func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	// validasi webhook token dari Xendit
	token := c.Get("x-callback-token")
	if token != h.webhookToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid webhook token",
		})
	}

	var req dto.WebhookRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid webhook payload")
	}

	if err := h.paymentService.HandleWebhook(req); err != nil {
		return response.InternalError(c, "Failed to process webhook: "+err.Error())
	}

	return response.OK(c, "Webhook processed", nil)
}
