package wildlifenl

import (
	"sync"
	"time"
)

// rateLimiter implements a per-token sliding-window rate limiter with
// escalating penalties. It allows bursts up to maxReq within window, then
// rejects excess requests with an exponentially growing cooldown. The penalty
// decays after the caller stops violating the limit for the decay duration.
//
// Stale bucket entries are removed by a background cleanup goroutine so the
// in-memory state does not grow without bound.
type rateLimiter struct {
	mu        sync.RWMutex
	buckets   map[string]*rateLimitState
	window    time.Duration
	maxReq    int
	baseDelay time.Duration
	maxDelay  time.Duration
	decay     time.Duration
}

// rateLimitState holds the mutable state for a single caller token.
type rateLimitState struct {
	mu            sync.Mutex
	requests      []time.Time // timestamps of requests within the sliding window
	penaltyLevel  int
	lastViolation time.Time
}

func newRateLimiter(window time.Duration, maxReq int, baseDelay, maxDelay, decay time.Duration) *rateLimiter {
	rl := &rateLimiter{
		buckets:   make(map[string]*rateLimitState),
		window:    window,
		maxReq:    maxReq,
		baseDelay: baseDelay,
		maxDelay:  maxDelay,
		decay:     decay,
	}
	go rl.cleanupLoop()
	return rl
}

// cleanupLoop periodically removes stale bucket entries.
func (rl *rateLimiter) cleanupLoop() {
	interval := rl.window
	if interval <= 0 {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup removes bucket entries that have no recent requests and whose
// penalty has fully decayed.
func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)
	for token, state := range rl.buckets {
		state.mu.Lock()
		state.requests = dropExpired(state.requests, cutoff)
		rl.decayPenalty(state, now)
		if len(state.requests) == 0 && state.penaltyLevel == 0 {
			delete(rl.buckets, token)
		}
		state.mu.Unlock()
	}
}

// dropExpired removes timestamps that fall on or before cutoff and returns
// the remaining slice. The input slice must be sorted in chronological order.
func dropExpired(requests []time.Time, cutoff time.Time) []time.Time {
	for i, t := range requests {
		if t.After(cutoff) {
			return requests[i:]
		}
	}
	return requests[:0]
}

// allow reports whether a request from token is permitted. If the request is
// rejected, ok is false and retryAfter indicates how long the caller must wait
// before trying again. All updates to the caller's state are atomic with
// respect to that caller.
func (rl *rateLimiter) allow(token string) (ok bool, retryAfter time.Duration) {
	state := rl.stateFor(token)
	state.mu.Lock()
	defer state.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Drop request timestamps that have fallen outside the sliding window.
	state.requests = dropExpired(state.requests, cutoff)

	// Decay penalty after a period of good behaviour.
	rl.decayPenalty(state, now)

	// If an escalation cooldown is active, reject immediately without
	// incrementing the violation count again.
	if !state.lastViolation.IsZero() {
		cooldown := rl.cooldown(state.penaltyLevel)
		expires := state.lastViolation.Add(cooldown)
		if now.Before(expires) {
			return false, expires.Sub(now)
		}
	}

	if len(state.requests) < rl.maxReq {
		state.requests = append(state.requests, now)
		return true, 0
	}

	// This request exceeds the limit: escalate and reject.
	state.penaltyLevel++
	state.lastViolation = now
	return false, rl.cooldown(state.penaltyLevel)
}

// decayPenalty reduces the penalty level based on how many full decay
// intervals have passed since the last violation. When the level reaches
// zero the violation timestamp is cleared.
func (rl *rateLimiter) decayPenalty(state *rateLimitState, now time.Time) {
	if state.penaltyLevel <= 0 || state.lastViolation.IsZero() {
		return
	}
	intervals := int(now.Sub(state.lastViolation) / rl.decay)
	if intervals <= 0 {
		return
	}
	state.penaltyLevel -= intervals
	if state.penaltyLevel <= 0 {
		state.penaltyLevel = 0
		state.lastViolation = time.Time{}
	}
}

// cooldown returns the penalty cooldown for the given violation level,
// capped at maxDelay. Level 1 yields baseDelay; each additional level
// doubles the delay until maxDelay is reached.
func (rl *rateLimiter) cooldown(level int) time.Duration {
	if level <= 0 {
		return 0
	}
	d := rl.baseDelay
	for i := 1; i < level; i++ {
		d *= 2
		if d >= rl.maxDelay {
			return rl.maxDelay
		}
	}
	return d
}

// stateFor returns the mutable state for token, creating it if necessary.
func (rl *rateLimiter) stateFor(token string) *rateLimitState {
	rl.mu.RLock()
	state, ok := rl.buckets[token]
	rl.mu.RUnlock()
	if ok {
		return state
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()
	state, ok = rl.buckets[token]
	if ok {
		return state
	}
	state = &rateLimitState{}
	rl.buckets[token] = state
	return state
}
