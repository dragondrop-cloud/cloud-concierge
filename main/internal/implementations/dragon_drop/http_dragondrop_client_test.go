package dragonDrop

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostLog(t *testing.T) {
	// Given
	dragonDropFactory := new(Factory)

	ctx := context.Background()
	mux := http.NewServeMux()

	mux.HandleFunc(
		"/log/",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		})

	server := httptest.NewServer(mux)
	defer server.Close()

	config := HTTPDragonDropClientConfig{
		APIPath:  server.URL,
		JobID:    "123",
		OrgToken: "123",
	}

	dragonDrop, err := dragonDropFactory.Instantiate("", config)
	assert.Nil(t, err)
	assert.NotNil(t, dragonDrop)

	// When
	err = dragonDrop.(*HTTPDragonDropClient).postLog(ctx, "Example log", false)
	assert.Nil(t, err)
}
