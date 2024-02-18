package vcs

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// Factory is a struct that generates implementations of interfaces.VCS
type Factory struct{}

// Instantiate returns an implementation of interfaces.VCS depending on the passed
// environment specification.
func (f *Factory) Instantiate(ctx context.Context, environment string, config Config, vcsSystem string) (interfaces.VCS, error) {
	switch environment {
	case "isolated":
		return new(IsolatedVCS), nil
	default:
		return f.bootstrappedVCS(ctx, config, vcsSystem)
	}
}

// bootstrappedVCS creates a complete implementation of the interfaces.VCS interface with
// configuration specified via environment variables.
func (f *Factory) bootstrappedVCS(ctx context.Context, config Config, vcsSystem string) (interfaces.VCS, error) {
	switch vcsSystem {
	case "github":
		return NewGitHub(config), nil
	default:
		log.Errorf("currently only GitHub is supported as a VCS option. %v was specified", vcsSystem)
		return nil, fmt.Errorf("currently only GitHub is supported as a VCS option. %v was specified", vcsSystem)
	}
}
