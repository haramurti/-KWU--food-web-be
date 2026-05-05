package dto

type CreateInvoiceRequest struct {
	ExternalID  string  `json:"external_id"`
	Amount      float64 `json:"amount"`
	PayerEmail  string  `json:"payer_email"`
	Description string  `json:"description"`
}

type CreateInvoiceResponse struct {
	InvoiceID  string  `json:"invoice_id"`
	InvoiceURL string  `json:"invoice_url"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

type WebhookRequest struct {
	ID         string  `json:"id"`
	ExternalID string  `json:"external_id"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
	PayerEmail string  `json:"payer_email"`
	PaidAt     string  `json:"paid_at"`
}
