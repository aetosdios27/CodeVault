package http

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimitedClient struct {
	limiter *rate.Limiter
	client  *http.Client
}

var (
	cfLimiter = rate.NewLimiter(rate.Every(500*time.Millisecond), 2) // 2 req/s
	lcLimiter = rate.NewLimiter(rate.Every(2*time.Second), 1)        // 0.5 req/s
	mu        sync.Mutex
)

func DoRequestWithRateLimit(req *http.Request, judge string) (*http.Response, error) {
	var limiter *rate.Limiter

	switch judge {
	case "codeforces":
		limiter = cfLimiter
	case "leetcode":
		limiter = lcLimiter
	default:
		limiter = rate.NewLimiter(rate.Inf, 0)
	}

	if err := limiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
