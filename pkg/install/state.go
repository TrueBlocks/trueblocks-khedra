package install

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// StepOrder defines the linear installation steps used for derivation and UI progress.
// Exported so other packages (templates / progress rendering) may rely on a single source.
var StepOrder = []string{"welcome", "paths", "chains", "index", "services", "logging", "summary"}

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

func currentStep() string { return StepOrder[0] }

func Handler(session *SessionStore, version string, configured bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sid, _ := session.Get()
		st := State{Configured: configured, CurrentStep: currentStep(), SessionID: sid, Version: version, Schema: 1}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(st)
	}
}

// EnsureID checks if the session already has an Id and, if yes, returns it, otherwise, it generates a new random Id.
func (session *SessionStore) EnsureID() string {
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.id != "" {
		return session.id
	}
	session.id = RandString(16)
	session.last = time.Now()
	return session.id

}

func RandString(n int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Enforce checks the posted session against the current session, handling inactivity and takeover logic.
// Returns status ("ok" or "conflict"), takeover (bool), current session id, and last activity time.
func (s *SessionStore) Enforce(postedSession string, inactivityWindow time.Duration) (status string, takeover bool, current string, last time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	current = s.id
	last = s.last
	now := time.Now()
	// If no session is set, accept postedSession as new session
	if current == "" {
		s.id = postedSession
		s.last = now
		return "ok", true, postedSession, now
	}
	// If posted session matches current, update last activity
	if postedSession == current {
		// Only update last if within inactivity window
		if now.Sub(s.last) < inactivityWindow {
			s.last = now
			return "ok", false, current, s.last
		}
		// If inactivity window passed, allow takeover
		return "ok", false, current, s.last
	}
	// If posted session is different
	if now.Sub(s.last) >= inactivityWindow {
		// Allow takeover
		s.id = postedSession
		s.last = now
		return "ok", true, postedSession, now
	}
	// Otherwise, conflict
	return "conflict", false, current, s.last
}
