package httptools

import (
	"fmt"
	"net/http"
	"time"
)

func NewHealthCheckHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("%s, I'm good at %v", name, time.Now())))
	}
}
