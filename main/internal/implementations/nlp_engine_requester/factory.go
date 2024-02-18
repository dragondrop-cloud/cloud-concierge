package nlpenginerequestor

import (
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct for creating different implementations of the NLPEngine interface.
type Factory struct{}

// Instantiate creates an implementation of the NLPEngine interface.
func (f *Factory) Instantiate(config HTTPNLPEngineClientConfig) (interfaces.NLPEngine, error) {
	return f.bootstrappedNLPEngine(config)
}

// bootstrappedNLPEngine instantiates an instance of a HTTPNLPEngineClient with the proper environment
// variables read in.
func (f *Factory) bootstrappedNLPEngine(config HTTPNLPEngineClientConfig) (interfaces.NLPEngine, error) {
	return NewHTTPNLPEngineClient(config), nil
}
