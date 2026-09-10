package handler

import (
	"analytics/internal/kafka"
	"analytics/internal/model"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	producer *kafka.Producer
}

func NewHandler(producer *kafka.Producer) *Handler {
return &Handler{
	producer: producer,
}
}

type EventRequest struct {
	EventType string 	`json:"event_type"`
  	Element string		`json:"element"`
  	Timestamp time.Time `json:"timestamp"`
  	Page string 		`json:"page"`
  	SessionID string 	`json:"session_id"` // uuid
  	Data map[string]any `json:"data"` // как проверять, тоже не должна быть пустой?
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	log.Println("received event request")
	ctx := r.Context()

	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "wrong JSON", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err !=nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := req.ToModel()
	log.Printf("полученные данные: %+v", event)

	if err := h.producer.Produce(ctx, event, "user-events"); err != nil {
		http.Error(w, "failed to publish event", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusAccepted)

}

func (e EventRequest) Validate() error { // как то не красив овыглядит
	if strings.TrimSpace(e.EventType) == "" {
        return errors.New("EventType is required")
    }
	if strings.TrimSpace(e.Element) == "" {
        return errors.New("Element is required")
    }
	if e.Timestamp.IsZero() {
    	return errors.New("Timestamp is required")
	}
	if e.Timestamp.After(time.Now()) {
		return errors.New("Timestamp cannot be in the future")
	}
	if strings.TrimSpace(e.Page) == "" {
        return errors.New("Page is required")
    }
	if strings.TrimSpace(e.SessionID) == "" {
        return errors.New("Session_id is required") // uuid checks?
    }
	if e.Data == nil { // точно ли все приходят с датой ?
        return errors.New("Data is required") 
    }
	return nil
}

func (e EventRequest) ToModel() model.Event {
    return model.Event{
    	EventType: e.EventType,
    	Element:   e.Element,
    	Timestamp: e.Timestamp,
    	Page:      e.Page,
    	SessionID: e.SessionID,
    	Data:      e.Data,
	}
}