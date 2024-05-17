package limiter

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/berachain/offchain-sdk/tools/store"
)

type Limiter struct {
	Store  store.Store
	Rate   int
	Period time.Duration
}

type Config struct {
	Enabled          bool
	TTL              time.Duration
	Rate             int
	Kind             string
	RedisAddr        string
	RedisClusterMode bool
	ProxyCount       int
}

func New(config Config) *Limiter {
	var lstore store.Store
	switch config.Kind {
	case "redis":
		lstore = store.NewRedisStore(config.TTL, config.RedisAddr, config.RedisClusterMode)
	case "memory":
		fallthrough
	default:
		lstore = store.NewInMemoryStore(config.TTL)
	}

	return &Limiter{
		Store:  lstore,
		Rate:   config.Rate,
		Period: time.Second,
	}
}

func Middleware(
	limiter *Limiter,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := getClientIPFromRequest(1, r)
			count, _, err := limiter.Store.Increment(r.Context(), key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if count > int64(limiter.Rate) {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getClientIPFromRequest(proxyCount int, r *http.Request) string {
	if proxyCount > 0 {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		if xForwardedFor != "" {
			xForwardedForParts := strings.Split(xForwardedFor, ",")
			// Avoid reading the user's forged request header by configuring the count of reverse proxies
			partIndex := len(xForwardedForParts) - proxyCount
			if partIndex < 0 {
				partIndex = 0
			}
			return strings.TrimSpace(xForwardedForParts[partIndex])
		}
	}

	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
