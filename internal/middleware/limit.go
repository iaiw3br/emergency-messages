package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	limit       = 1
	burst       = 3
	maxLastSeen = 3 * time.Minute
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

// getVisitor
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(limit, burst)
		visitors[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupVisitors clean from the cache
func cleanupVisitors() {
	mu.Lock()
	for ip, v := range visitors {
		if time.Since(v.lastSeen) > maxLastSeen {
			delete(visitors, ip)
		}
	}
	mu.Unlock()
}

// LimitRequests rate limiter by IP
func LimitRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		cleanupVisitors()
		limiter := getVisitor(ip)
		if limiter.Allow() == false {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
