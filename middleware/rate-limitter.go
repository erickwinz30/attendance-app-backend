package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter menyimpan informasi jumlah request per IP
type RateLimiter struct {
	requests map[string][]time.Time // IP -> daftar waktu request
	mu       sync.Mutex             // untuk thread-safe
	limit    int                    // jumlah request maksimal
	window   time.Duration          // jendela waktu (1 menit)
}

// NewRateLimiter membuat rate limiter baru
// limit: jumlah request per window (contoh: 60)
// window: durasi jendela waktu (contoh: 1 menit)
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow mengecek apakah request diizinkan untuk IP tertentu
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Jika IP belum ada, buat entry baru
	if _, exists := rl.requests[ip]; !exists {
		rl.requests[ip] = []time.Time{}
	}

	// Hapus request yang lebih lama dari window
	var validRequests []time.Time
	for _, reqTime := range rl.requests[ip] {
		if now.Sub(reqTime) < rl.window {
			validRequests = append(validRequests, reqTime)
		}
	}
	rl.requests[ip] = validRequests

	// Cek apakah sudah mencapai limit
	if len(rl.requests[ip]) >= rl.limit {
		return false // Request ditolak
	}

	// Tambahkan request baru
	rl.requests[ip] = append(rl.requests[ip], now)
	return true // Request diizinkan
}

// RateLimitMiddleware adalah middleware untuk menerapkan rate limiter
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ambil IP address dari request
			ip := r.RemoteAddr

			// Cek apakah request diizinkan
			if !limiter.Allow(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests) // HTTP 429
				w.Write([]byte(`{"error":"Too many requests. Max 60 requests per minute"}`))
				return
			}

			// Lanjutkan ke handler berikutnya
			next.ServeHTTP(w, r)
		})
	}
}
