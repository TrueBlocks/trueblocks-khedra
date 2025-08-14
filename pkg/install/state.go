package install

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var steps = []string{"welcome", "paths", "index", "chains", "services", "logging", "summary"}

type State struct {
	Configured  bool   `json:"configured"`
	CurrentStep string `json:"currentStep"`
	SessionID   string `json:"sessionId"`
	Version     string `json:"version"`
	Schema      int    `json:"schema"`
}

type SessionStore struct {
	mu   sync.Mutex
	id   string
	last time.Time
}

func NewSessionStore() *SessionStore { return &SessionStore{} }

func (s *SessionStore) Touch(id string) {
	s.mu.Lock()
	s.id = id
	s.last = time.Now()
	s.mu.Unlock()
}

func (s *SessionStore) Get() (string, time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.id, s.last
}

func currentStep() string { return steps[0] }

func Handler(session *SessionStore, version string, configured bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, _ := session.Get()
		st := State{Configured: configured, CurrentStep: currentStep(), SessionID: sid, Version: version, Schema: 1}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(st)
	}
}
