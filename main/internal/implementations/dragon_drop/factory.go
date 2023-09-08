package dragondrop

import (
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct for creating different implementations of the DragonDrop interface.
type Factory struct {
}

// Instantiate creates an implementation of the DragonDrop interface.
func (f *Factory) Instantiate(environment string, config HTTPDragonDropClientConfig) (interfaces.DragonDrop, error) {
	switch environment {
	case "isolated":
		return new(IsolatedDragonDrop), nil
	default:
		return f.bootstrappedDragonDrop(config)
	}
}

// bootstrappedDragonDrop instantiates an instance of a HTTPDragonDropClient with the proper environment
// variables read in.
func (f *Factory) bootstrappedDragonDrop(config HTTPDragonDropClientConfig) (interfaces.DragonDrop, error) {
	return NewHTTPDragonDropClient(config), nil
}
