package schedulessrv

import (
	"bytes"
	"complete-api/internal/core/domain"
	"complete-api/internal/core/ports"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type service struct {
	redisRepo ports.RedisRepository
}

func New(redisRepo ports.RedisRepository) *service {
	return &service{
		redisRepo: redisRepo,
	}
}

func (s *service) Run() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.processDueMessages()
	}
}

func parseNextCronTime(expr string, from int64) (int64, error) {
	sched, err := cron.ParseStandard(expr)
	if err != nil {
		return 0, err
	}
	next := sched.Next(time.Unix(from, 0))
	return next.Unix(), nil
}

func (s *service) processDueMessages() {
	now := time.Now().Unix()

	users, err := s.redisRepo.GetAllScheduledUsers()
	if err != nil {
		log.Println("Erro ao buscar usuários:", err)
		return
	}

	for _, userID := range users {
		// Busca IDs dos agendamentos vencidos
		ids, err := s.redisRepo.GetZRangeByScore(userID, 0, now)
		if err != nil {
			log.Printf("Erro ao buscar mensagens para %s: %v", userID, err)
			continue
		}

		for _, schedule := range ids {
			id := schedule.ID // já vem montado no struct

			// Envia a mensagem
			if err := s.sendMessage(schedule); err != nil {
				log.Println("Erro ao enviar mensagem:", err)
				continue
			}

			// Aplica recorrência
			if schedule.CronExpr != "" {
				// Gera próximo horário
				nextTime, err := parseNextCronTime(schedule.CronExpr, now)
				if err != nil {
					log.Println("Erro ao calcular próxima recorrência:", err)
					continue
				}

				// Condição 1: Se tiver limite por quantidade
				if schedule.Repeats > 0 {
					schedule.Repeats--
					if schedule.Repeats == 0 {
						// Remove completamente
						s.redisRepo.RemoveZMember(userID, id)
						continue
					}
				}

				// Condição 2: Se tiver limite por data
				if schedule.Until > 0 && nextTime > schedule.Until {
					s.redisRepo.RemoveZMember(userID, id)
					continue
				}

				// Atualiza `SendAt` para o próximo agendamento
				schedule.SendAt = nextTime
				if err := s.redisRepo.UpdateScheduledMessage(userID, id, schedule); err != nil {
					log.Println("Erro ao reagendar mensagem:", err)
					continue
				}

			} else {
				// Não é recorrente: deleta
				s.redisRepo.RemoveZMember(userID, id)
			}
		}
	}
}

func (s *service) sendMessage(msg domain.ScheduleMessage) error {
	// Monta a URL usando o canal
	apiURL := fmt.Sprintf("https://api.l2msg.com/whatsapp/message/sendText/%s", msg.Channel)
	apiKey := "429683C4C977415CAAFCCE10F7D57E11" // Substitua pela sua chave de API

	// Monta o payload esperado pela API
	payload := map[string]string{
		"number": msg.Phone,
		"text":   msg.Text,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Cria a requisição POST
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	// Define headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)

	// Executa a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Valida o status de resposta
	if resp.StatusCode >= 300 {
		return fmt.Errorf("erro ao enviar mensagem: status %d", resp.StatusCode)
	}

	return nil
}

func (s *service) GetSchedules(username string) ([]domain.ScheduleMessage, error) {
	schedules, err := s.redisRepo.GetZRangeByScore(username, 0, 9999999999)
	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return []domain.ScheduleMessage{}, nil
	}

	return schedules, nil
}

func (s *service) CreateSchedule(username string, schedule domain.ScheduleMessage) (domain.ScheduleMessage, error) {
	schedule.ID = uuid.NewString()

	err := s.redisRepo.AddScheduledMessage(username, schedule)
	if err != nil {
		return domain.ScheduleMessage{}, err
	}

	return schedule, nil
}

func (s *service) UpdateSchedule(username string, scheduleID string, schedule domain.ScheduleMessage) (domain.ScheduleMessage, error) {
	// Garante que o ID seja consistente com o parâmetro
	schedule.ID = scheduleID

	// Atualiza o schedule no Redis
	err := s.redisRepo.UpdateScheduledMessage(username, scheduleID, schedule)
	if err != nil {
		return domain.ScheduleMessage{}, err
	}

	return schedule, nil
}

func (s *service) DeleteSchedule(username string, scheduleID string) error {
	return s.redisRepo.RemoveZMember(username, scheduleID)
}
