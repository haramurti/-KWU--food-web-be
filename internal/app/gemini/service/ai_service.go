package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/haramurti/-KWU--food-web-be/internal/app/gemini/dto"
)

const geminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"

const systemPrompt = `Kamu adalah asisten AI untuk brand health bar bernama "NutriBar". 
Tugasmu adalah membantu pengguna menemukan varian NutriBar yang paling cocok untuk mereka.

Varian yang tersedia:

1. NutriBar BOOST — 280 kkal, 12g protein. Cocok untuk atlet & pre-workout. Bahan: oat, madu, kurma, almond.
2. NutriBar LEAN  — 180 kkal, 20g protein. Cocok untuk diet & turun berat badan. Bahan: whey isolate, chia, dark choco.
3. NutriBar GLOW  — 210 kkal, 8g protein. Cocok untuk kulit sehat & hormonal. Bahan: collagen, mixed berry, flaxseed.
4. NutriBar CALM  — 200 kkal, 6g protein. Cocok untuk stres & susah tidur. Bahan: ashwagandha, lavender, pisang.
5. NutriBar KID   — 160 kkal, 5g protein. Cocok untuk anak 5–12 tahun. Bahan: oat, madu, pisang, susu, no preservatives.

Panduan:
- Tanya dulu kebutuhan pengguna sebelum merekomendasikan
- Rekomendasikan 1–2 varian dengan alasan jelas
- Bahasa menyesuaikan user (Indonesia/Inggris)
- Singkat, maksimal 3–4 paragraf`

// ── internal Gemini API structs ──────────────────────────────────────────────

type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"system_instruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	GenerationConfig  generationCfg   `json:"generationConfig"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type generationCfg struct {
	MaxOutputTokens int     `json:"maxOutputTokens"`
	Temperature     float64 `json:"temperature"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type GeminiService struct {
	apiKey string
	client *http.Client
}

func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (g *GeminiService) Chat(history []dto.ChatMessageDTO, userMessage string) (string, error) {
	contents := make([]geminiContent, 0, len(history)+1)

	for _, msg := range history {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, geminiContent{
			Role:  role,
			Parts: []geminiPart{{Text: msg.Content}},
		})
	}

	contents = append(contents, geminiContent{
		Role:  "user",
		Parts: []geminiPart{{Text: userMessage}},
	})

	reqBody := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: systemPrompt}},
		},
		Contents: contents,
		GenerationConfig: generationCfg{
			MaxOutputTokens: 512,
			Temperature:     0.7,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	url := fmt.Sprintf("%s?key=%s", geminiEndpoint, g.apiKey)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body error: %w", err)
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(raw, &gemResp); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	if gemResp.Error != nil {
		return "", fmt.Errorf("gemini API error: %s", gemResp.Error.Message)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from gemini")
	}

	return gemResp.Candidates[0].Content.Parts[0].Text, nil
}
