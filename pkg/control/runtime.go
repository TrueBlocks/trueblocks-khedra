package control

import "sync"

// RuntimeState tracks ephemeral pause state without mutating persisted config.
type RuntimeState struct {
	mu     sync.RWMutex
	paused map[string]bool
}

func NewRuntimeState() *RuntimeState { return &RuntimeState{paused: map[string]bool{}} }

func (r *RuntimeState) Pause(names ...string) (changed []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, n := range names {
		if !r.paused[n] {
			r.paused[n] = true
			changed = append(changed, n)
		}
	}
	return
}

func (r *RuntimeState) Unpause(names ...string) (changed []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, n := range names {
		if r.paused[n] {
			delete(r.paused, n)
			changed = append(changed, n)
		}
	}
	return
}

func (r *RuntimeState) IsPaused(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.paused[name]
}

func (r *RuntimeState) Snapshot() map[string]bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp := make(map[string]bool, len(r.paused))
	for k, v := range r.paused {
		cp[k] = v
	}
	return cp
}
