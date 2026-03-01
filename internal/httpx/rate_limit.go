package httpx

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keithics/devops-dashboard/api/internal/httperr"
	"golang.org/x/time/rate"
)

type ipVisitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*ipVisitor
	limit    rate.Limit
	burst    int
	ttl      time.Duration
}

func newIPRateLimiter(limit rate.Limit, burst int, ttl time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		visitors: make(map[string]*ipVisitor),
		limit:    limit,
		burst:    burst,
		ttl:      ttl,
	}
}

func (r *ipRateLimiter) allow(ip string) bool {
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	for key, visitor := range r.visitors {
		if now.Sub(visitor.lastSeen) > r.ttl {
			delete(r.visitors, key)
		}
	}

	visitor, ok := r.visitors[ip]
	if !ok {
		visitor = &ipVisitor{
			limiter:  rate.NewLimiter(r.limit, r.burst),
			lastSeen: now,
		}
		r.visitors[ip] = visitor
	}

	visitor.lastSeen = now
	return visitor.limiter.Allow()
}

// RateLimitByIP returns a Gin middleware that limits requests by client IP.
func RateLimitByIP(limit rate.Limit, burst int, ttl time.Duration) gin.HandlerFunc {
	limiter := newIPRateLimiter(limit, burst, ttl)

	return func(c *gin.Context) {
		if limiter.allow(c.ClientIP()) {
			c.Next()
			return
		}

		httperr.Abort(c, httperr.New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "rate limit exceeded"))
	}
}
