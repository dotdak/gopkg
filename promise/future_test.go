package promise

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFutureToReturnChannels(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(w, "Hello World")
	}))
	defer ts.Close()

	f := New(func(context.Context) (interface{}, error) {
		resp, err := http.Get(ts.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	})

	futureChan := Future(context.Background(), f)

	expectedChan := &Result{Data: []byte("Hello World"), Err: nil}

	assert.Equal(t, expectedChan, <-futureChan)
}
