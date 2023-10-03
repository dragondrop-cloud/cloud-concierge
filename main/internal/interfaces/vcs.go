package interfaces

import "github.com/stretchr/testify/mock"

// VCS interface for interacting with version control systems. Since all major VCS systems
// are git-based (GitHub, GitLab, BitBucket, etc.), they can share the same interface.
type VCS interface {

	// Clone pulls a remote repository's contents into local memory.
	Clone() error

	// AddChanges adds all code changes to be included in the next commit.
	AddChanges() error

	// Checkout creates a new branch within the remote repository.
	Checkout(jobName string) error

	// Commit commits code changes to the current branch of the remote repository.
	Commit() error

	// Push pushes current branch to remote repository.
	Push() error

	// OpenPullRequest opens a new pull request of committed changes to the remote repository,
	// and returns the url of this pull request
	OpenPullRequest(jobName string) (string, error)

	// GetID returns a string which is a random, 10 character unique identifier
	// for a dragondrop built commit/pull request
	GetID() (string, error)

	SetToken(token string)
}

// VCSMock implements the VCS interface solely for testing purposes.
type VCSMock struct {
	mock.Mock
}

// Clone pulls a remote repository's contents into local memory.
func (m *VCSMock) Clone() error {
	args := m.Called()
	return args.Error(0)
}

// AddChanges adds all code changes to be included in the next commit.
func (m *VCSMock) AddChanges() error {
	args := m.Called()
	return args.Error(0)
}

// Checkout creates a new branch within the remote repository.
func (m *VCSMock) Checkout(_ string) error {
	args := m.Called()
	return args.Error(0)
}

// Commit commits code changes to the current branch of the remote repository.
func (m *VCSMock) Commit() error {
	args := m.Called()
	return args.Error(0)
}

// Push pushes current branch to remote repository.
func (m *VCSMock) Push() error {
	args := m.Called()
	return args.Error(0)
}

// OpenPullRequest opens a new pull request of committed changes to the remote repository,
// and returns the url of this pull request
func (m *VCSMock) OpenPullRequest(_ string) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// GetID returns a string which is a random, 10 character unique identifier
// for a dragondrop built commit/pull request
func (m *VCSMock) GetID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// SetToken sets the token for the VCS
func (m *VCSMock) SetToken(token string) {
	m.Called(token)
}
