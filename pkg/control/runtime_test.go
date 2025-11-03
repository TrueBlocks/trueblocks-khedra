package control

import (
	"math/rand"
	"sync"
	"testing"
)

func TestRuntimeStateBasicBehavior(t *testing.T) {
	t.Run("Pause idempotent", func(t *testing.T) {
		r := NewRuntimeState()
		changed := r.Pause("a", "a")
		if len(changed) != 1 || !r.IsPaused("a") {
			t.Fatalf("expected single change and paused=true, got %v", changed)
		}
		unchanged := r.Pause("a")
		if len(unchanged) != 0 {
			t.Fatalf("expected no changes on idempotent pause")
		}
	})

	t.Run("Unpause idempotent", func(t *testing.T) {
		r := NewRuntimeState()
		r.Pause("a")
		changed := r.Unpause("a", "a")
		if len(changed) != 1 || r.IsPaused("a") {
			t.Fatalf("expected unpause, got %v paused=%v", changed, r.IsPaused("a"))
		}
		if snap := r.Snapshot(); len(snap) != 0 {
			t.Fatalf("expected empty snapshot after unpause, got %v", snap)
		}
	})

	t.Run("Multiple names", func(t *testing.T) {
		r := NewRuntimeState()
		r.Pause("a", "b")
		r.Unpause("b")
		snap := r.Snapshot()
		if !snap["a"] || snap["b"] {
			t.Fatalf("expected a paused, b removed; got %v", snap)
		}
	})

	t.Run("IsPaused reflects state", func(t *testing.T) {
		r := NewRuntimeState()
		r.Pause("svc")
		if !r.IsPaused("svc") || r.IsPaused("other") {
			t.Fatalf("IsPaused not reflecting state")
		}
	})

	t.Run("Snapshot isolation", func(t *testing.T) {
		r := NewRuntimeState()
		r.Pause("a")
		snap := r.Snapshot()
		r.Unpause("a")
		if !snap["a"] {
			t.Fatalf("expected original snapshot to retain value")
		}
		if r.IsPaused("a") {
			t.Fatalf("expected runtime state unpaused now")
		}
	})
}

func TestRuntimeStateConcurrencyAndRaces(t *testing.T) {
	t.Run("Parallel pause/unpause", func(t *testing.T) {
		r := NewRuntimeState()
		names := []string{"a", "b", "c", "d"}
		var wg sync.WaitGroup
		rand.Seed(1)
		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 200; j++ {
					n := names[rand.Intn(len(names))]
					if rand.Intn(2) == 0 {
						r.Pause(n)
					} else {
						r.Unpause(n)
					}
					if j%10 == 0 {
						_ = r.Snapshot()
					}
				}
			}()
		}
		wg.Wait()
		snap := r.Snapshot()
		for k := range snap {
			found := false
			for _, n := range names {
				if k == n {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("unexpected key in snapshot: %s", k)
			}
		}
	})

	t.Run("High churn snapshot", func(t *testing.T) {
		r := NewRuntimeState()
		names := []string{"a", "b", "c", "d"}
		var wg sync.WaitGroup
		rand.Seed(2)
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					n := names[rand.Intn(len(names))]
					if rand.Intn(2) == 0 {
						r.Pause(n)
					} else {
						r.Unpause(n)
					}
				}
			}()
		}
		for i := 0; i < 20; i++ {
			_ = r.Snapshot()
		}
		wg.Wait()
		snap := r.Snapshot()
		for k := range snap {
			found := false
			for _, n := range names {
				if k == n {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("unexpected key in snapshot: %s", k)
			}
		}
	})
}
