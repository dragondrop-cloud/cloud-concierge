package vcs

// Config contains the values that allow for authentication and the specific repo
// traits needed.
type Config struct {

	// VCSRepo is the full path of the repo containing a customer's infrastructure specification.
	VCSRepo string `required:"true"`

	// PullReviewers is the name of the pull request reviewer who will be tagged on the opened pull request.
	PullReviewers []string `default:"NoReviewer"`
}
