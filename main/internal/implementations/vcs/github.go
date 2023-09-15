package vcs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v45/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/dragondrop-cloud/cloud-concierge/main/internal/interfaces"
)

// GitHub struct implements the VCS interface.
type GitHub struct {
	// authBasic is the authentication information needed to perform generic git operations via
	authBasic *http.BasicAuth

	// config contains the values that allow for authentication and the specific repo
	// traits needed.
	config Config

	// defaultBranch is the name of the default branch of the repository.
	defaultBranch string

	// dragonDrop is needed to inform cloned status
	dragonDrop interfaces.DragonDrop

	// ID is a string which is a random, 10 character unique identifier
	// for a cloud-concierge built commit/pull request
	ID string

	// newBranchName is the name of the new branch name for the new pull request.
	newBranchName string

	// oauth2Client is an authenticated client that is able
	// to access the customer's GitHub account. Primarily used for opening pull requests.
	oauth2Client *github.Client

	// repository is a code repository object from the go-git package which represents the customer's
	// code repository containing IaC.
	repository *git.Repository

	// workTree is the working tree object which references repository
	workTree *git.Worktree
}

// NewGitHub creates a new instance of the GitHub struct.
func NewGitHub(ctx context.Context, dragonDrop interfaces.DragonDrop, config Config) interfaces.VCS {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.VCSToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	authenticatedClient := github.NewClient(tc)
	dragonDrop.PostLog(ctx, "Created VCS client.")

	return &GitHub{
		config: config,
		authBasic: &http.BasicAuth{
			Username: config.VCSUser,
			Password: config.VCSToken,
		},
		oauth2Client: authenticatedClient,
		dragonDrop:   dragonDrop,
	}
}

// GetDefaultBranch returns the default branch of the repository.
func (g *GitHub) GetDefaultBranch() error {
	vcsCloneURLArray := strings.Split(g.config.VCSRepo, "/")
	repoName := strings.Replace(vcsCloneURLArray[len(vcsCloneURLArray)-1], ".git", "", -1)
	repoOwner := vcsCloneURLArray[len(vcsCloneURLArray)-2]

	repoReference, _, err := g.oauth2Client.Repositories.Get(context.Background(), repoOwner, repoName)
	if err != nil {
		return fmt.Errorf("[g.oauth2Client.Repositories.Get][%v]", err)
	}

	g.defaultBranch = repoReference.GetDefaultBranch()

	logrus.Debugf("[Github] Default branch is %v", g.defaultBranch)
	return nil
}

// GetID returns a string which is a random, 10 character unique identifier
// for a cloud-concierge built commit/pull request
func (g *GitHub) GetID() (string, error) {
	if strings.Trim(g.ID, "") == "" {
		return "", errors.New("[vcs][get_id][id not generated]")
	}

	logrus.Debugf("[Github] ID is %v", g.ID)
	return g.ID, nil
}

// Clone pulls a remote repository's contents into local memory.
func (g *GitHub) Clone() error {
	cloneOptions := &git.CloneOptions{
		Auth:     g.authBasic,
		URL:      g.config.VCSRepo,
		Progress: os.Stdout,
	}

	// Cleaning out the existing repository folder. Cannot clone into an already existing directory.
	err := os.RemoveAll("./repo/")
	if err != nil {
		return err
	}

	repo, err := git.PlainClone("./repo/", false, cloneOptions)
	if err != nil {
		return err
	}

	g.repository = repo

	logrus.Debugf("[Github] Cloned repo %v", g.config.VCSRepo)
	return nil
}

// AddChanges adds all code changes to be included in the next commit.
func (g *GitHub) AddChanges() error {
	logrus.Debugf("[Github] Adding changes to repo %v", g.config.VCSRepo)
	addOptions := &git.AddOptions{
		All: true,
	}

	err := g.workTree.AddWithOptions(addOptions)

	if err != nil {
		return fmt.Errorf("[vcs][add_changed][error in worktree.AddWithOptions]%w", err)
	}

	return nil
}

// Checkout creates a new branch within the remote repository.
func (g *GitHub) Checkout(jobName string) error {
	lowerJobName := strings.ToLower(jobName)
	jobNameSplit := strings.Split(lowerJobName, " ")
	cleanJobName := strings.Join(jobNameSplit, "_")

	branchUniqueID := time.Now().Format("2006-01-02-15-04")

	newBranchName := fmt.Sprintf(
		"feature/cloud_concierge_%v_%v",
		cleanJobName,
		branchUniqueID,
	)

	g.newBranchName = newBranchName

	branchName := plumbing.NewBranchReferenceName(newBranchName)

	checkoutOptions := &git.CheckoutOptions{
		Branch: branchName,
		Create: true,
	}

	workTree, err := g.repository.Worktree()

	if err != nil {
		return fmt.Errorf("[vcs][checkout][error in creating worktree]%w", err)
	}

	err = workTree.Checkout(checkoutOptions)

	if err != nil {
		return fmt.Errorf("[vcs][checkout][error in checking out a new branch for the suggested changes]%w", err)
	}

	g.workTree = workTree
	g.ID = branchUniqueID

	logrus.Debugf("[Github] Checked out branch %v", g.newBranchName)
	return nil
}

// Commit commits code changes to the current branch of the remote repository.
func (g *GitHub) Commit() error {
	logrus.Debugf("[Github] Committing changes to repo %v", g.config.VCSRepo)

	commitOptions := &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "dragondrop.cloud",
			Email: "cloud-concierge@dragondrop.cloud",
			When:  time.Now(),
		},
	}

	commitHash, err := g.workTree.Commit("build: cloud-concierge results", commitOptions)

	if err != nil {
		return fmt.Errorf("[vcs][commit][error in worktree.AddWithOptions]%w", err)
	}

	fmt.Printf("Commit made with hash: %v\n", commitHash)

	return nil
}

// Push pushes current branch to remote repository.
func (g *GitHub) Push() error {
	logrus.Debugf("[Github] Pushing changes to repo %v", g.config.VCSRepo)

	pushOptions := &git.PushOptions{
		Auth:     g.authBasic,
		Progress: os.Stdout,
	}

	err := g.repository.Push(pushOptions)

	if err != nil {
		return fmt.Errorf("[vcs][push][error in repository.Push]%w", err)
	}

	return nil
}

// OpenPullRequest opens a new pull request of committed changes to the remote repository.
func (g *GitHub) OpenPullRequest(jobName string) (string, error) {
	prTitle := fmt.Sprintf("%v - %v", jobName, g.ID)
	logrus.Debugf("[Github] Opening PR with title %v", prTitle)

	reportContent, err := os.ReadFile("state_of_cloud/report.md")
	if err != nil {
		return "", fmt.Errorf("error in loading state of cloud report: %v", err)
	}

	prComment := string(reportContent)

	err = g.GetDefaultBranch()
	if err != nil {
		return "", fmt.Errorf("[g.GetDefaultBranch]%v", err)
	}

	newPR := &github.NewPullRequest{
		Title:               &prTitle,
		Head:                &g.newBranchName,
		Base:                &g.defaultBranch,
		Body:                &prComment,
		MaintainerCanModify: github.Bool(true),
	}

	orgName, repoName, err := g.extractOrgAndRepoName(g.config.VCSRepo)

	if err != nil {
		return "", fmt.Errorf("[extractOrgAndRepoName] %v", err)
	}

	pr, _, err := g.oauth2Client.PullRequests.Create(
		context.Background(),
		orgName,
		repoName,
		newPR,
	)

	if err != nil {
		return "", fmt.Errorf("error in github.PullRequests.Create(): %v", err)
	}

	if g.config.PullReviewers[0] != "NoReviewer" {
		rr := github.ReviewersRequest{
			Reviewers: g.config.PullReviewers,
		}

		_, _, err = g.oauth2Client.PullRequests.RequestReviewers(
			context.Background(),
			orgName,
			repoName,
			pr.GetNumber(),
			rr,
		)

		if err != nil {
			return "", fmt.Errorf("error in github.PullRequests.RequestReviewers(): %v", err)
		}
	}

	logrus.Debugf("[Github] PR opened with url %v", pr.GetURL())
	return pr.GetURL(), nil
}

// extractOrgAndRepoName pulls out the organization and repository name from the
// repositories full path.
func (g *GitHub) extractOrgAndRepoName(repoFullPath string) (string, string, error) {
	r, err := regexp.Compile(`github.com/(.*?)/(.*?).git$`)

	if err != nil {
		return "", "", fmt.Errorf("[extract_org_and_repo_name][error in regexp.Compile]%w", err)
	}

	org := r.FindStringSubmatch(repoFullPath)[1]

	repo := r.FindStringSubmatch(repoFullPath)[2]

	return org, repo, nil
}
