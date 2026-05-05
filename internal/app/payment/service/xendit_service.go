package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/haramurti/-KWU--food-web-be/internal/app/payment/contract"
	"github.com/haramurti/-KWU--food-web-be/internal/app/payment/dto"
)

var _ contract.PaymentService = (*XenditService)(nil)

const xenditBaseURL = "https://api.xendit.co"

type XenditService struct {
	apiKey       string
	webhookToken string
	client       *http.Client
}

// ── xendit internal request/response structs ─────────────────────────────────

type xenditInvoiceRequest struct {
	ExternalID     string   `json:"external_id"`
	Amount         float64  `json:"amount"`
	PayerEmail     string   `json:"payer_email"`
	Description    string   `json:"description"`
	Currency       string   `json:"currency"`
	PaymentMethods []string `json:"payment_methods,omitempty"`
}

type xenditInvoiceResponse struct {
	ID         string  `json:"id"`
	ExternalID string  `json:"external_id"`
	InvoiceURL string  `json:"invoice_url"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	// error fields
	ErrorCode string `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
}

// ── constructor ───────────────────────────────────────────────────────────────

func NewXenditService(apiKey, webhookToken string) *XenditService {
	return &XenditService{
		apiKey:       apiKey,
		webhookToken: webhookToken,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (x *XenditService) authHeader() string {
	// Xendit pakai basic auth: apiKey + ":" di-encode base64
	raw := x.apiKey + ":"
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(raw))
}

func (x *XenditService) doRequest(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal error: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, xenditBaseURL+path, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("build request error: %w", err)
	}

	req.Header.Set("Authorization", x.authHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := x.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read body error: %w", err)
	}

	return raw, resp.StatusCode, nil
}

// ── CreateInvoice ─────────────────────────────────────────────────────────────

func (x *XenditService) CreateInvoice(req dto.CreateInvoiceRequest) (*dto.CreateInvoiceResponse, error) {
	xenditReq := xenditInvoiceRequest{
		ExternalID:  req.ExternalID,
		Amount:      req.Amount,
		PayerEmail:  req.PayerEmail,
		Description: req.Description,
		Currency:    "IDR",
		PaymentMethods: []string{
			"BCA", "BNI", "BRI", "MANDIRI",
			"OVO", "DANA", "LINKAJA",
			"QRIS",
		},
	}

	raw, statusCode, err := x.doRequest(http.MethodPost, "/v2/invoices", xenditReq)
	if err != nil {
		return nil, err
	}

	var xenditResp xenditInvoiceResponse
	if err := json.Unmarshal(raw, &xenditResp); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		return nil, fmt.Errorf("xendit error [%d]: %s — %s",
			statusCode, xenditResp.ErrorCode, xenditResp.Message)
	}

	return &dto.CreateInvoiceResponse{
		InvoiceID:  xenditResp.ID,
		InvoiceURL: xenditResp.InvoiceURL,
		Amount:     xenditResp.Amount,
		Status:     xenditResp.Status,
	}, nil
}

// ── HandleWebhook ─────────────────────────────────────────────────────────────

func (x *XenditService) HandleWebhook(req dto.WebhookRequest) error {
	// di sini nanti bisa tambahin logic:
	// - update status order di DB
	// - kirim notif ke user
	// - dll
	switch req.Status {
	case "PAID":
		// TODO: mark order as paid
		fmt.Printf("[WEBHOOK] Order %s PAID — amount: %.0f\n", req.ExternalID, req.Amount)
	case "EXPIRED":
		// TODO: mark order as expired
		fmt.Printf("[WEBHOOK] Order %s EXPIRED\n", req.ExternalID)
	default:
		fmt.Printf("[WEBHOOK] Order %s status: %s\n", req.ExternalID, req.Status)
	}

	return nil
}
