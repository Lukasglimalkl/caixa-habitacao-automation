package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/queue"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

func main() {

	// Inicializa o logger
	logger.Init()

	workerID := os.Getenv("WORKER_ID")
	if workerID == "" {
		workerID = "worker-1"
	}

	logger.Info(fmt.Sprintf("🚀 Worker %s iniciando...", workerID))

	// Conecta na fila Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	q := queue.NewRedisQueue(redisAddr)
	logger.Info(fmt.Sprintf("✅ Conectado ao Redis: %s", redisAddr))

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Loop principal do worker
	go func() {
		for {
			logger.Info(fmt.Sprintf("[%s] 🔍 Buscando próximo job...", workerID))

			// Pega próximo job da fila
			job, err := q.GetNextJob()
			if err != nil {
				logger.Error(fmt.Sprintf("[%s] ❌ Erro ao buscar job: %v", workerID, err))
				time.Sleep(5 * time.Second)
				continue
			}

			if job == nil {
				// Fila vazia, aguarda
				time.Sleep(2 * time.Second)
				continue
			}

			logger.Info(fmt.Sprintf("[%s] 📋 Processando job %s (CPF: %s)", workerID, job.ID, job.CPF))

			// Atualiza status para processing
			job.Status = "processing"
			q.UpdateJob(job)

			// Executa automação
			bot := automation.NewCaixaBot(true)
			
			req := models.LoginAndSearchRequest{
				Username: job.Username,
				Password: job.Password,
				CPF:      job.CPF,
			}

			response, err := bot.LoginAndSearch(req.Username, req.Password, req.CPF)
			
			if err != nil {
				logger.Error(fmt.Sprintf("[%s] ❌ Erro no job %s: %v", workerID, job.ID, err))
				q.FailJob(job.ID, err.Error())
				continue
			}

			// Serializa resultado
			resultJSON, _ := json.Marshal(response.Data)

			// Marca como completo
			q.CompleteJob(job.ID, string(resultJSON))

			logger.Info(fmt.Sprintf("[%s] ✅ Job %s completado!", workerID, job.ID))
		}
	}()

	// Aguarda sinal de stop
	<-stop
	logger.Info(fmt.Sprintf("[%s] 🛑 Worker parando...", workerID))
}