package model
import "time"

type Event struct {
    EventType string         `json:"event_type"`
    Element   string         `json:"element"`
    Timestamp time.Time      `json:"timestamp"`
    Page      string         `json:"page"`
    SessionID string         `json:"session_id"` // uuid
    Data      map[string]any `json:"data,omitempty"` // может быть пустым?
}