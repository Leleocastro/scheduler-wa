package api_gateway

import (
	"bytes"
	"complete-api/internal/core/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
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
	apiKeys, err := getAPIKey(s.baseURL, username)
	if err != nil {
		return err
	}

	if len(apiKeys.Data) > 0 {
		fmt.Println("Chave de API já existe!")
		return nil
	}

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
	response, err := getAPIKey(s.baseURL, username)
	if err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", fmt.Errorf("erro: chave de API não encontrada para o consumidor no Kong")
	}

	sort.Slice(response.Data, func(i, j int) bool {
		return response.Data[i].CreatedAt > response.Data[j].CreatedAt
	})

	fmt.Println("Chave de API obtida com sucesso!")

	return response.Data[0].Key, nil
}

func (s *kongAPI) RemoveRateLimitConsumer(username, route string) error {
	url := fmt.Sprintf("%s/consumers/%s/plugins", s.baseURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro: status %s ao buscar plugins do consumidor no Kong", resp.Status)
	}

	var response domain.PluginsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("erro ao decodificar resposta do Kong: %w", err)
	}

	for _, plugin := range response.Data {
		url := fmt.Sprintf("%s/routes/%s", s.baseURL, plugin.Route.ID)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("erro: status %s ao buscar rota no Kong", resp.Status)
		}

		var routeData domain.Route

		if err := json.NewDecoder(resp.Body).Decode(&routeData); err != nil {
			return fmt.Errorf("erro ao decodificar resposta do Kong: %w", err)
		}

		log.Println("name: ", routeData.Name)

		if *routeData.Name == route {
			url = fmt.Sprintf("%s/consumers/%s/plugins/%s", s.baseURL, username, plugin.ID)

			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				return fmt.Errorf("erro: status %s ao remover plugin do consumidor no Kong", resp.Status)
			}

			fmt.Println("Limite de taxa removido com sucesso!")
			return nil
		}
	}

	return fmt.Errorf("erro: plugin de limite de taxa não encontrado para o consumidor no Kong")
}

func getAPIKey(baseURL, username string) (*domain.ApiKeyResponse, error) {
	url := fmt.Sprintf("%s/consumers/%s/key-auth", baseURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro: status %s ao buscar chave de API para o consumidor no Kong", resp.Status)
	}

	var response domain.ApiKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta do Kong: %w", err)
	}

	return &response, nil
}

func (s *kongAPI) RemoveACL(username, group string) error {
	url := fmt.Sprintf("%s/consumers/%s/acls", s.baseURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro: status %s ao buscar ACLs do consumidor no Kong", resp.Status)
	}

	var response domain.ACLsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("erro ao decodificar resposta do Kong: %w", err)
	}

	group = strings.ToLower(group)

	for _, acl := range response.Data {
		if acl.Group == group {
			url = fmt.Sprintf("%s/consumers/%s/acls/%s", s.baseURL, username, acl.ID)

			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				return fmt.Errorf("erro ao criar requisição HTTP: %w", err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("erro ao fazer requisição para o Kong: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusNoContent {
				return fmt.Errorf("erro: status %s ao remover ACL do consumidor no Kong", resp.Status)
			}

			fmt.Println("ACL removido com sucesso!")
			return nil
		}
	}

	return fmt.Errorf("erro: ACL não encontrada para o consumidor no Kong")
}
