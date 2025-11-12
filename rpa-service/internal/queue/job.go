package queue

import (
	"encoding/json"
	"time"
)

// Job - representa um trabalho na fila
type Job struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CPF       string    `json:"cpf"`
	Status    string    `json:"status"` // pending, processing, completed, failed
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToJSON - converte Job para JSON
func (j *Job) ToJSON() (string, error) {
	data, err := json.Marshal(j)
	return string(data), err
}

// FromJSON - converte JSON para Job
func FromJSON(data string) (*Job, error) {
	var job Job
	err := json.Unmarshal([]byte(data), &job)
	return &job, err
}