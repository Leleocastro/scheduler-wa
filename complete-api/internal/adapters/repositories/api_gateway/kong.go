package api_gateway

import (
	"bytes"
	"complete-api/internal/core/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type kongAPI struct {
	baseURL string
}

func New(baseURL string) *kongAPI {
	return &kongAPI{
		baseURL: baseURL,
	}
}

func (s *kongAPI) CreateConsumer(username, customID string) error {
	// Definindo a URL do endpoint de consumidores
	url := fmt.Sprintf("%s/consumers", s.baseURL)

	// Criando a estrutura do consumidor
	consumer := domain.Consumer{
		Username: username,
		CustomID: customID,
	}

	// Convertendo a estrutura para JSON
	payload, err := json.Marshal(consumer)
	if err != nil {
		return fmt.Errorf("erro ao criar payload do consumidor: %w", err)
	}

	// Fazendo a requisição POST para o Kong
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("KONG_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erro: status %s ao criar consumidor no Kong", resp.Status)
	}

	fmt.Println("Consumidor criado com sucesso!")
	return nil
}

func (s *kongAPI) RateLimitConsumer(username, route string, rateLimit int) error {
	// Definindo a URL do endpoint de limites de taxa
	url := fmt.Sprintf("%s/consumers/%s/plugins", s.baseURL, username)

	fmt.Println("URL:", url)

	// Criando a estrutura do limite de taxa
	rateLimiting := map[string]any{
		"name":  "rate-limiting",
		"route": route,
		"config": map[string]int{
			"day": rateLimit,
		},
	}

	// Convertendo a estrutura para JSON
	payload, err := json.Marshal(rateLimiting)
	if err != nil {
		return fmt.Errorf("erro ao criar payload do limite de taxa: %w", err)
	}

	// Fazendo a requisição POST para o Kong
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", os.Getenv("KONG_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erro: status %s ao criar limite de taxa para o consumidor no Kong", resp.Status)
	}

	fmt.Println("Limite de taxa criado com sucesso!")
	return nil
}

func (s *kongAPI) CreateJWTFirebaseConsumer(username string) error {
	url := fmt.Sprintf("%s/consumers/%s/plugins", s.baseURL, username)

	// Criando a estrutura do plugin
	jwtFirebase := map[string]any{
		"name": "jwt-firebase",
		"config": map[string]string{
			"project_id":       "scheduler-wa",
			"jwt_service_user": username,
		},
	}

	// Convertendo a estrutura para JSON
	payload, err := json.Marshal(jwtFirebase)
	if err != nil {
		return fmt.Errorf("erro ao criar payload do plugin: %w", err)
	}

	// Fazendo a requisição POST para o Kong
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erro: status %s ao criar plugin para o consumidor no Kong", resp.Status)
	}

	fmt.Println("Plugin criado com sucesso!")

	return nil
}
