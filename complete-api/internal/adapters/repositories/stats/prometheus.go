package stats

import (
	"complete-api/internal/core/domain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type prometheusAPI struct {
	baseURL string
}

func New(baseURL string) *prometheusAPI {
	return &prometheusAPI{
		baseURL: baseURL,
	}
}

func (p *prometheusAPI) GetUsageByConsumer(username string, startDate, endDate int64) (domain.UsageResponse, error) {
	query := fmt.Sprintf(`round(sum(increase(kong_http_requests_total{consumer="%s"}[1d])) by (consumer, service))`, username)
	escapedQuery := url.QueryEscape(query)

	url := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=1d", p.baseURL, escapedQuery, startDate, endDate)

	fmt.Println("Querying Prometheus...")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return domain.UsageResponse{}, fmt.Errorf("error creating HTTP request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return domain.UsageResponse{}, fmt.Errorf("error making request to Prometheus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Erro ao ler o corpo da resposta: %v", err)
		} else {
			bodyString := string(bodyBytes)
			log.Println("ERROR Body: ", bodyString)
		}
		return domain.UsageResponse{}, fmt.Errorf("error: status %s querying Prometheus", resp.Status)
	}

	var response domain.UsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return domain.UsageResponse{}, fmt.Errorf("error decoding response body: %w", err)
	}

	fmt.Println("Query successful!")
	return response, nil

}
