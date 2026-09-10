package main

import (
	"analytics/internal/kafka"
	"analytics/internal/store"
	"analytics/internal/usecase"
	"context"
	
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("OTMENA")
		cancel()
	}()

	err := godotenv.Load()
	if err != nil {
        log.Println("Файл .env не найден, используются только переменные окружения")
    }
	
	databaseURL := os.Getenv("DATABASE_URL")
	brokers := os.Getenv("KAFKA_BROKERS")

	db, err := store.NewPostgres(databaseURL)
	if err != nil {
		log.Fatal("ошибка подключения к БД: ", err)
	}
	defer db.Close()

	uc:= usecase.NewAnalyticsUsecase(db)

	conshandler := kafka.NewConsumerHandler(uc)

	consumer, err := kafka.NewConsumer([]string{brokers}, "worker-group")
	if err != nil {
		log.Fatal("ошибка создания consumer group: ", err)
	}
	defer consumer.Close()

	if err := consumer.Consume(ctx, []string{"user-events"}, conshandler); err != nil && ctx.Err() == nil {
		log.Fatal("ошибка consumer: ", err)
	}

	log.Println("Worker остановлен")
}