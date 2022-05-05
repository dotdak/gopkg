package httptools

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type HttpRequest[R, S comparable] func(ctx context.Context, req *R) (S, error)

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(err.Error()))
}

const CtxHTTPHeader = "http_header"
const CtxHTTPHost = "http_host"

func HandlerWrap[R, S comparable](f HttpRequest[R, S]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		inBuf, err := io.ReadAll(r.Body)
		if err != nil {
			BadRequest(w, err)
			return
		}
		var req R
		if len(inBuf) > 0 {
			if err := json.Unmarshal(inBuf, &req); err != nil {
				BadRequest(w, err)
				return
			}
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxHTTPHeader, r.Header)
		ctx = context.WithValue(ctx, CtxHTTPHost, r.Host)
		val, err := f(ctx, &req)
		if err != nil {
			BadRequest(w, err)
			return
		}
		buf, err := json.Marshal(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(buf)
		return
	}
}

func GetHost(ctx context.Context) string {
	val, ok := ctx.Value(CtxHTTPHost).(string)
	if ok {
		return val
	}

	return ""
}
