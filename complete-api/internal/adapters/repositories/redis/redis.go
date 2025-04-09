package redis

import (
	"complete-api/internal/core/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func New(addr string, password string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,     // ex: "localhost:6379"
		Password: password, // vazio se não tiver senha
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Erro ao conectar no Redis: %v", err))
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisClient) AddScheduledMessage(userID string, schedule domain.ScheduleMessage) error {
	key := fmt.Sprintf("scheduled:messages:%s", userID)
	dataKey := fmt.Sprintf("scheduled:data:%s", schedule.ID)

	payload, err := json.Marshal(schedule)
	if err != nil {
		return err
	}

	pipe := r.client.TxPipeline()

	// 1. Armazena o ID no ZSet
	pipe.ZAdd(r.ctx, key, redis.Z{
		Score:  float64(schedule.SendAt),
		Member: schedule.ID,
	})

	// 2. Armazena os dados reais do schedule separados
	pipe.Set(r.ctx, dataKey, payload, 0)

	// 3. Garante que o usuário esteja no índice de usuários
	pipe.SAdd(r.ctx, "scheduled:users", userID)

	_, execErr := pipe.Exec(r.ctx)
	return execErr
}

func (r *RedisClient) GetZRangeByScore(userID string, min, max int64) ([]domain.ScheduleMessage, error) {
	key := fmt.Sprintf("scheduled:messages:%s", userID)

	// 1. Pega os IDs no range
	ids, err := r.client.ZRangeByScore(r.ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", min),
		Max: fmt.Sprintf("%d", max),
	}).Result()
	if err != nil {
		return nil, err
	}

	var schedules []domain.ScheduleMessage
	for _, id := range ids {
		dataKey := fmt.Sprintf("scheduled:data:%s", id)

		// 2. Busca os dados com base no ID
		data, err := r.client.Get(r.ctx, dataKey).Result()
		if err != nil {
			// Se o dado não existir, remove o ID do ZSet e continua
			if err == redis.Nil {
				r.client.ZRem(r.ctx, key, id)
				continue
			}
			return nil, err
		}

		// 3. Converte para struct
		var msg domain.ScheduleMessage
		if err := json.Unmarshal([]byte(data), &msg); err != nil {
			continue // Ignora mensagens mal formatadas
		}

		schedules = append(schedules, msg)
	}

	return schedules, nil
}

func (r *RedisClient) RemoveZMember(userID string, scheduleID string) error {
	key := fmt.Sprintf("scheduled:messages:%s", userID)
	dataKey := fmt.Sprintf("scheduled:data:%s", scheduleID)

	pipe := r.client.TxPipeline()

	// 1. Remove do ZSet
	pipe.ZRem(r.ctx, key, scheduleID)

	// 2. Remove o conteúdo salvo
	pipe.Del(r.ctx, dataKey)

	// 3. Se for o último agendamento do usuário, remove do índice
	zcardCmd := pipe.ZCard(r.ctx, key)

	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return err
	}

	// Verifica o número de agendamentos restantes
	count, err := zcardCmd.Result()
	if err != nil {
		return err
	}
	if count == 0 {
		if err := r.client.SRem(r.ctx, "scheduled:users", userID).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisClient) UpdateScheduledMessage(userID, id string, msg domain.ScheduleMessage) error {
	dataKey := fmt.Sprintf("scheduled:data:%s", id)
	zsetKey := fmt.Sprintf("scheduled:messages:%s", userID)

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	pipe := r.client.TxPipeline()
	pipe.Set(r.ctx, dataKey, payload, 0)
	pipe.ZAdd(r.ctx, zsetKey, redis.Z{
		Score:  float64(msg.SendAt),
		Member: id,
	})
	_, err = pipe.Exec(r.ctx)
	return err
}

func (r *RedisClient) GetAllScheduledUsers() ([]string, error) {
	return r.client.SMembers(r.ctx, "scheduled:users").Result()
}
