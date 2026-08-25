package main

import (
	"analytics/internal/kafka"
	handler "analytics/internal/transport/http"
	"log"
	"net/http"
)

func main() {
	p, err := kafka.NewProducer("localhost:9092")
	if err != nil {
		log.Fatalf("Не удалось создать кафку: %v", err)
	}
	
	defer p.Close()

	h := handler.NewHandler(p)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/events", h.CreateEvent)

	if err := http.ListenAndServe(":8080", handler.CORSMiddleware(mux)); err != nil {
   		log.Fatal(err)
	}
}