package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

const (
	currentVersion = "v0.1.11"
	footer         = "###########################################################################################################################################################"
)

// VersionChecker is a struct that contains the http client used to make http requests to GitHub and validate the latest version of cloud-concierge.
type VersionChecker struct {
	// httpClient is a http client shared across all http requests within this package.
	httpClient http.Client
}

// checkCloudConciergeVersion checks the current version of cloud-concierge against the latest stable version.
func (v *VersionChecker) checkCloudConciergeVersion() {
	latestStableVersion, err := v.getLatestStableCloudConciergeVersion()
	if err != nil {
		log.Errorf("[version_checker] Error checking latest cloud concierge latestStableVersion: %s", err.Error())
		return
	}

	cloudConciergeTitle := fmt.Sprintf("################################################################# Cloud Concierge %s ##################################################################", currentVersion)

	message, level := v.getCloudConciergeVersionMessage(currentVersion, latestStableVersion)

	if level == "red" {
		color.Red(cloudConciergeTitle)
		color.Red(fmt.Sprintf("# %s #", message))
		color.Red(footer)
	} else {
		color.Green(cloudConciergeTitle)
		color.Green(fmt.Sprintf("############################### %s ##################################", message))
		color.Green(footer)
	}
}

// GithubRelease is a struct that contains the tag_name of a github release.
type GithubRelease struct {
	TagName string `json:"tag_name"`
}

// getLatestStableCloudConciergeVersion gets the latest stable version of cloud-concierge from github.
func (v *VersionChecker) getLatestStableCloudConciergeVersion() (string, error) {
	log.Debugf("Checking for latest latestStableVersion of cloud-concierge...")
	response, err := v.httpClient.Get("https://api.github.com/repos/dragondrop-cloud/cloud-concierge/releases")
	if err != nil {
		return "", fmt.Errorf("[version_checker] error getting latest cloud concierge latestStableVersion: %s", err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[version_checker] error getting latest cloud concierge latestStableVersion: %s", response.Status)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("[version_checker] error reading response body: %s", err.Error())
	}

	var githubReleases []GithubRelease
	err = json.Unmarshal(body, &githubReleases)
	if err != nil {
		return "", fmt.Errorf("[version_checker] error unmarshalling response body: %s", err.Error())
	}

	for _, release := range githubReleases {
		if !strings.Contains(release.TagName, "-beta") {
			return release.TagName, nil
		}
	}

	return "", nil
}

// getCloudConciergeVersionMessage returns a message and level based on the current and latest version of cloud-concierge.
func (v *VersionChecker) getCloudConciergeVersionMessage(currentVersion, latestVersion string) (string, string) {
	if latestVersion == "" {
		return fmt.Sprintf("You are currently running version %s of cloud-concierge, this is the latest version.", currentVersion), "green"
	}

	if latestVersion != currentVersion {
		return fmt.Sprintf("You are currently running version %s of cloud-concierge, the latest version is %s, run docker pull dragondrop/cloud-concierge:latest to update.", currentVersion, latestVersion), "red"
	}

	return fmt.Sprintf("You are currently running version %s of cloud-concierge, this is the latest version.", currentVersion), "green"
}
