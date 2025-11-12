package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	JobsQueue       = "rpa:jobs:queue"
	JobsProcessing  = "rpa:jobs:processing"
	JobsCompleted   = "rpa:jobs:completed"
	JobsKeyPrefix   = "rpa:job:"
)

type RedisQueue struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisQueue - cria nova fila Redis
func NewRedisQueue(addr string) *RedisQueue {
	client := redis.NewClient(&redis.Options{
		Addr: addr, // Ex: "localhost:6379"
	})

	return &RedisQueue{
		client: client,
		ctx:    context.Background(),
	}
}

// AddJob - adiciona job na fila
func (q *RedisQueue) AddJob(username, password, cpf string) (string, error) {
	job := &Job{
		ID:        uuid.New().String(),
		Username:  username,
		Password:  password,
		CPF:       cpf,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Salva job no Redis
	jobJSON, err := job.ToJSON()
	if err != nil {
		return "", err
	}

	jobKey := fmt.Sprintf("%s%s", JobsKeyPrefix, job.ID)
	err = q.client.Set(q.ctx, jobKey, jobJSON, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	// Adiciona na fila
	err = q.client.RPush(q.ctx, JobsQueue, job.ID).Err()
	if err != nil {
		return "", err
	}

	return job.ID, nil
}

// GetNextJob - pega próximo job da fila
func (q *RedisQueue) GetNextJob() (*Job, error) {
	// Move job da fila para processing (atomic)
	result, err := q.client.BLMove(q.ctx, JobsQueue, JobsProcessing, "LEFT", "RIGHT", 5*time.Second).Result()
	
	if err == redis.Nil {
		return nil, nil // Fila vazia
	}
	if err != nil {
		return nil, err
	}

	// Busca dados do job
	jobKey := fmt.Sprintf("%s%s", JobsKeyPrefix, result)
	jobJSON, err := q.client.Get(q.ctx, jobKey).Result()
	if err != nil {
		return nil, err
	}

	return FromJSON(jobJSON)
}

// UpdateJob - atualiza status do job
func (q *RedisQueue) UpdateJob(job *Job) error {
	job.UpdatedAt = time.Now()
	
	jobJSON, err := job.ToJSON()
	if err != nil {
		return err
	}

	jobKey := fmt.Sprintf("%s%s", JobsKeyPrefix, job.ID)
	return q.client.Set(q.ctx, jobKey, jobJSON, 24*time.Hour).Err()
}

// CompleteJob - marca job como completo
func (q *RedisQueue) CompleteJob(jobID string, result string) error {
	jobKey := fmt.Sprintf("%s%s", JobsKeyPrefix, jobID)
	jobJSON, err := q.client.Get(q.ctx, jobKey).Result()
	if err != nil {
		return err
	}

	job, err := FromJSON(jobJSON)
	if err != nil {
		return err
	}

	job.Status = "completed"
	job.Result = result

	// Atualiza job
	if err := q.UpdateJob(job); err != nil {
		return err
	}

	// Remove de processing, adiciona em completed
	q.client.LRem(q.ctx, JobsProcessing, 1, jobID)
	q.client.RPush(q.ctx, JobsCompleted, jobID)

	return nil
}

// FailJob - marca job como falho
func (q *RedisQueue) FailJob(jobID string, errorMsg string) error {
	jobKey := fmt.Sprintf("%s%s", JobsKeyPrefix, jobID)
	jobJSON, err := q.client.Get(q.ctx, jobKey).Result()
	if err != nil {
		return err
	}

	job, err := FromJSON(jobJSON)
	if err != nil {
		return err
	}

	job.Status = "failed"
	job.Error = errorMsg

	// Atualiza job
	if err := q.UpdateJob(job); err != nil {
		return err
	}

	// Remove de processing
	q.client.LRem(q.ctx, JobsProcessing, 1, jobID)

	return nil
}

// GetJobStatus - busca status de um job
func (q *RedisQueue) GetJobStatus(jobID string) (*Job, error) {
	jobKey := fmt.Sprintf("%s%s", JobsKeyPrefix, jobID)
	jobJSON, err := q.client.Get(q.ctx, jobKey).Result()
	if err != nil {
		return nil, err
	}

	return FromJSON(jobJSON)
}