package main

import (
	"analytics/internal/kafka"
	handler "analytics/internal/transport/http"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	brokers := os.Getenv("KAFKA_BROKERS")

	p, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatalf("Не удалось создать кафку: %v", err)
	}
	defer p.Close()

	h := handler.NewHandler(p)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/events", h.CreateEvent)
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	server := &http.Server{Addr: ":8080",
							Handler: handler.CORSMiddleware(mux),
						}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("стартовал на :8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("ошибка HTTP server: %v", err)
		}
	}()

	<-sigChan
	log.Println("получен сигнал завершения")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("ошибка graceful shutdown HTTP server: %v", err)
	}

	log.Println("API остановлен")
}