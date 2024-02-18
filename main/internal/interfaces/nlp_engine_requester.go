package interfaces

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// NLPEngine is the interface for communicating with the external NLPEngine API
type NLPEngine interface {
	// PostNLPEngine sends a request to the NLPEngine endpoint
	// and saves results into local container memory.
	PostNLPEngine(ctx context.Context) error
}

// NLPEngineMock is a struct that implements the NLPEngine interface solely for the purpose
// of testing with the testify library.
type NLPEngineMock struct {
	mock.Mock
}
