package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	maxAttempts   = 5
	windowSeconds = 15 * 60 // 15 minutos
)

type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
}

var loginLimiter = &rateLimiter{
	attempts: make(map[string][]time.Time),
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-time.Duration(windowSeconds) * time.Second)

	var recent []time.Time
	for _, t := range rl.attempts[key] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= maxAttempts {
		rl.attempts[key] = recent
		return false
	}

	recent = append(recent, now)
	rl.attempts[key] = recent
	return true
}

func RateLimitLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)

		if !loginLimiter.allow(ip) {
			http.Error(w, `{"error": "muitas tentativas de login, tente novamente em alguns minutos"}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}
