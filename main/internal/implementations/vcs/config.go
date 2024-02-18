package vcs

// Config contains the values that allow for authentication and the specific repo
// traits needed.
type Config struct {
	// VCSPat is the personal access token for the customer's VCS account. It must have
	// the necessary permissions to open pull requests and push commits to the VCSRepo specified below.
	VCSPat string `required:"true"`

	// VCSRepo is the full path of the repo containing a customer's infrastructure specification.
	VCSRepo string `required:"true"`

	// PullReviewers is the name of the pull request reviewer who will be tagged on the opened pull request.
	PullReviewers []string `default:"NoReviewer"`
}
