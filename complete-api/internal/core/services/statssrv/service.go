package statssrv

import (
	"complete-api/internal/core/domain"
	"complete-api/internal/core/ports"
	"fmt"
	"strconv"
	"time"
)

type service struct {
	prometheusRepo ports.StatsRepository
}

func New(prometheusRepo ports.StatsRepository) *service {
	return &service{
		prometheusRepo: prometheusRepo,
	}
}

func (s *service) GetUsageByConsumer(username, startDate, endDate string) (domain.Usage, error) {

	var startTimestamp, endTimestamp int64

	if startDate == "" {
		now := time.Now().UTC().Add(24 * time.Hour)
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		startTimestamp = startOfDay.Unix()
	} else {
		parsedStartDate, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return domain.Usage{}, fmt.Errorf("invalid start_date format")
		}
		startOfDay := time.Date(parsedStartDate.Year(), parsedStartDate.Month(), parsedStartDate.Day(), 0, 0, 0, 0, time.UTC)
		startTimestamp = startOfDay.Unix()
	}

	if endDate == "" {
		now := time.Now().UTC().Add(24 * time.Hour)
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
		endTimestamp = endOfDay.Unix()
	} else {
		parsedEndDate, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return domain.Usage{}, fmt.Errorf("invalid end_date format")
		}
		endOfDay := time.Date(parsedEndDate.Year(), parsedEndDate.Month(), parsedEndDate.Day(), 23, 59, 59, 0, time.UTC)
		endTimestamp = endOfDay.Unix()
	}

	resp, err := s.prometheusRepo.GetUsageByConsumer(username, startTimestamp, endTimestamp)
	if err != nil {
		return domain.Usage{}, err
	}

	if len(resp.Data.Result) == 0 {
		return domain.Usage{}, fmt.Errorf("no data found for the given consumer")
	}

	return mapUsageResponseToUsage(resp), nil
}

func mapUsageResponseToUsage(response domain.UsageResponse) domain.Usage {
	usageMap := make(map[string]map[string][]domain.UsageItem)

	// Itera sobre os resultados para estruturar os dados
	for _, result := range response.Data.Result {
		service := result.Metric.Service
		consumer := result.Metric.Consumer

		// Inicializa o mapa para o serviço se não existir
		if _, exists := usageMap[service]; !exists {
			usageMap[service] = make(map[string][]domain.UsageItem)
		}

		// Inicializa a lista para o consumidor se não existir
		if _, exists := usageMap[service][consumer]; !exists {
			usageMap[service][consumer] = []domain.UsageItem{}
		}

		// Adiciona os valores
		for _, value := range result.Values {
			if len(value) == 2 {
				timestamp, ok1 := value[0].(float64)
				countStr, ok2 := value[1].(string)

				if ok1 && ok2 {
					// Converte o timestamp para data
					date := time.Unix(int64(timestamp), 0).Add(-24 * time.Hour).Format("2006-01-02")

					// Converte o count para inteiro
					count, err := strconv.Atoi(countStr)
					if err == nil {
						usageMap[service][consumer] = append(usageMap[service][consumer], domain.UsageItem{
							Date:  date,
							Count: count,
						})
					}
				}
			}
		}
	}

	// Constrói o objeto final
	usage := domain.Usage{
		Usage: []domain.UsageService{},
	}

	for service, consumers := range usageMap {
		for _, items := range consumers {
			usage.Usage = append(usage.Usage, domain.UsageService{
				Service: service,
				Values:  items,
			})
		}
	}

	return usage
}
