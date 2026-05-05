package contract

import "github.com/haramurti/-KWU--food-web-be/internal/app/payment/dto"

type PaymentService interface {
	CreateInvoice(req dto.CreateInvoiceRequest) (*dto.CreateInvoiceResponse, error)
	HandleWebhook(req dto.WebhookRequest) error
}
