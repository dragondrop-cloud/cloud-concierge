package interfaces

import (
	"github.com/stretchr/testify/mock"
)

// DragonDrop is the interface for communicating with the external DragonDrop API
type DragonDrop interface {
}

// DragonDropMock is a struct that implements the DragonDrop interface solely for the purpose
// of testing with the testify library.
type DragonDropMock struct {
	mock.Mock
}
