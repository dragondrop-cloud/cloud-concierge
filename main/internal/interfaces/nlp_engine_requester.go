package interfaces

import (
	"github.com/stretchr/testify/mock"
)

// NLPEngine is the interface for communicating with the external NLPEngine API
type NLPEngine interface{}

// NLPEngineMock is a struct that implements the NLPEngine interface solely for the purpose
// of testing with the testify library.
type NLPEngineMock struct {
	mock.Mock
}
