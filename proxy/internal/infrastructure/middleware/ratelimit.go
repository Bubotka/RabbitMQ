package middleware

import (
	"fmt"
	"gitlab.com/ptflp/gopubsub/queue"
	"net/http"
	"time"
)

type RateLimit struct {
	limit           int
	currentRequests int
	queuer          queue.MessageQueuer
}

func NewRateLimit(limit int, queuer queue.MessageQueuer) *RateLimit {
	return &RateLimit{limit: limit, queuer: queuer}
}

func (rl *RateLimit) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rl.currentRequests >= rl.limit {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Println("Всё кончились попытки")

			err := rl.queuer.Publish("limit", []byte("Ну усё ты попал, нет больше попыток!"))
			if err != nil {
				fmt.Println("Чёт не получить отправить сообщение")
			}
			return
		}
		rl.currentRequests++
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimit) ResetCurrentRequestsPerTime(timeInSec int) {
	ticker := time.NewTicker(time.Duration(timeInSec) * time.Second)
	go func() {
		for range ticker.C {
			rl.currentRequests = 0
		}
	}()
}
