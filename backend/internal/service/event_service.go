package service

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type EventType string

const (
	EventStatusUpdate EventType = "status_update"
)

type Event struct {
	SubmissionID uuid.UUID   `json:"submission_id"`
	Type         EventType   `json:"type"`
	Status       string      `json:"status"`
	Message      string      `json:"message"`
}

type EventService struct {
	clients    map[uuid.UUID]chan Event
	clientsMu  sync.RWMutex
}

func NewEventService() *EventService {
	return &EventService{
		clients: make(map[uuid.UUID]chan Event),
	}
}

func (s *EventService) Subscribe(submissionID uuid.UUID) chan Event {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	ch := make(chan Event, 10)
	s.clients[submissionID] = ch
	return ch
}

func (s *EventService) Unsubscribe(submissionID uuid.UUID) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	if ch, ok := s.clients[submissionID]; ok {
		close(ch)
		delete(s.clients, submissionID)
	}
}

func (s *EventService) Broadcast(event Event) {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	if ch, ok := s.clients[event.SubmissionID]; ok {
		ch <- event
	}
}

func (s *EventService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// SSE implementation
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	submissionIDStr := r.URL.Query().Get("id")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		http.Error(w, "Invalid submission ID", http.StatusBadRequest)
		return
	}

	ch := s.Subscribe(submissionID)
	defer s.Unsubscribe(submissionID)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(w, "data: {\"connected\": true}\n\n")
	flusher.Flush()

	for {
		select {
		case event := <-ch:
			fmt.Fprintf(w, "event: %s\ndata: {\"status\": \"%s\", \"message\": \"%s\"}\n\n", event.Type, event.Status, event.Message)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
