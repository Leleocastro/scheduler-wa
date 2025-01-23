package api_gateway

import (
	"bytes"
	"complete-api/internal/core/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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
	url := fmt.Sprintf("%s/consumers", s.baseURL)

	consumer := domain.Consumer{
		Username: username,
		CustomID: customID,
	}

	payload, err := json.Marshal(consumer)
	if err != nil {
		return fmt.Errorf("erro ao criar payload do consumidor: %w", err)
	}

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
		return fmt.Errorf("erro: status %s ao criar consumidor no Kong", resp.Status)
	}

	fmt.Println("Consumidor criado com sucesso!")
	return nil
}

func (s *kongAPI) RateLimitConsumer(username, route string, rateLimit int) error {
	url := fmt.Sprintf("%s/consumers/%s/plugins", s.baseURL, username)

	rateLimiting := map[string]any{
		"name": "rate-limiting",
		"route": map[string]string{
			"name": route,
		},
		"config": map[string]int{
			"day": rateLimit,
		},
	}

	payload, err := json.Marshal(rateLimiting)
	if err != nil {
		return fmt.Errorf("erro ao criar payload do limite de taxa: %w", err)
	}

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
		return fmt.Errorf("erro: status %s ao criar limite de taxa para o consumidor no Kong", resp.Status)
	}

	fmt.Println("Limite de taxa criado com sucesso!")
	return nil
}

func (s *kongAPI) CreateACL(username, group string) error {
	url := fmt.Sprintf("%s/consumers/%s/acls", s.baseURL, username)

	groupData := map[string]any{
		"group": group,
	}

	payload, err := json.Marshal(groupData)
	if err != nil {
		return fmt.Errorf("erro ao criar payload do plugin: %w", err)
	}

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
		return fmt.Errorf("erro: status %s ao criar ACL para o consumidor no Kong", resp.Status)
	}

	fmt.Println("ACL criado com sucesso!")

	return nil
}

func (s *kongAPI) CreateAPIKey(username string) error {
	url := fmt.Sprintf("%s/consumers/%s/key-auth", s.baseURL, username)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("erro: status %s ao criar chave de API para o consumidor no Kong", resp.Status)
	}

	fmt.Println("Chave de API criada com sucesso!")

	return nil
}

func (s *kongAPI) GetAPIKey(username string) (string, error) {
	url := fmt.Sprintf("%s/consumers/%s/key-auth", s.baseURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro: status %s ao buscar chave de API para o consumidor no Kong", resp.Status)
	}

	var response domain.ApiKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta do Kong: %w", err)
	}

	sort.Slice(response.Data, func(i, j int) bool {
		return response.Data[i].CreatedAt > response.Data[j].CreatedAt
	})

	fmt.Println("Chave de API obtida com sucesso!")

	return response.Data[0].Key, nil
}
