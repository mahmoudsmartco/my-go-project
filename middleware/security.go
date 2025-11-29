package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// -------------------------
//   Password Hashing
// -------------------------

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// -------------------------
//        Rate Limit
// -------------------------

var visitors = make(map[string]time.Time)
var mu sync.Mutex

// Correct function signature
func RateLimit(next http.Handler, limit time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mu.Lock()
		last, exists := visitors[r.RemoteAddr]
		if exists && time.Since(last) < limit {
			mu.Unlock()
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		visitors[r.RemoteAddr] = time.Now()
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
