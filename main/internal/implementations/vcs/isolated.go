package vcs

import (
	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// IsolatedVCS is a struct that implements interfaces.VCS for
// the purpose of end-to-end testing.
type IsolatedVCS struct {
	CloneTimes uint
}

// NewIsolatedVCS returns a new instance of IsolatedVCS.
func NewIsolatedVCS(times uint) interfaces.VCS {
	return &IsolatedVCS{
		CloneTimes: times,
	}
}

// Clone pulls a remote repository's contents into local memory.
func (v *IsolatedVCS) Clone() error {
	v.CloneTimes++
	return nil
}

// AddChanges adds all code changes to be included in the next commit.
func (v *IsolatedVCS) AddChanges() error {
	return nil
}

// Checkout creates a new branch within the remote repository.
func (v *IsolatedVCS) Checkout(_ string) error {
	return nil
}

// Commit commits code changes to the current branch of the remote repository.
func (v *IsolatedVCS) Commit() error {
	return nil
}

// Push pushes current branch to remote repository.
func (v *IsolatedVCS) Push() error {
	return nil
}

// OpenPullRequest opens a new pull request of committed changes to the remote repository,
// and returns the url of this pull request
func (v *IsolatedVCS) OpenPullRequest(_ string) (string, error) {
	return "", nil
}

// GetID returns a string which is a random, 10 character unique identifier
// for a dragondrop built commit/pull request
func (v *IsolatedVCS) GetID() (string, error) {
	return "", nil
}

// SetToken sets the token for the VCS
func (v *IsolatedVCS) SetToken() {
}
