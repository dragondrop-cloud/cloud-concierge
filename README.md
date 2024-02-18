# cloud-concierge
<p align="center">
<img width="375" src=./images/cloud-concierge-logo.png>
</p>

<p align = "center">
<a href="https://goreportcard.com/report/github.com/dragondrop-cloud/cloud-concierge/main" alt="Go Report">
   <img src="https://img.shields.io/badge/Go_Report-A+-green" />
</a>

<a href="https://github.com/dragondrop-cloud/cloud-concierge/actions/workflows/ci.yml?query=branch%3Aprod" alt="Coverage Report">
   <img src="https://img.shields.io/badge/Tests-passing-darkgreen" />
</a>

<a href="https://hub.docker.com/r/dragondropcloud/cloud-concierge/tags" alt="Latest Docker Version">
   <img src="https://img.shields.io/badge/docker-v0.2.0-blue" />
</a>

<a href="https://hub.docker.com/r/dragondropcloud/cloud-concierge" alt="Total Downloads">
   <img src="https://img.shields.io/badge/downloads-13.6k-maroon" />
</a>

<h3 align="center">
<a href="https://github.com/dragondrop-cloud/cloud-concierge-example/pull/3" target="_blank">Example Output</a> |
<a href="https://docs.cloudconcierge.io" target="_blank">Docs</a>
</h3>

## Why cloud-concierge?
cloud-concierge is a container that integrates with your existing Terraform management set up.
All results and codified resources are output via a digestible [Pull Request](https://github.com/dragondrop-cloud/cloud-concierge-example/pull/3) to a repository of your choice, providing you with a "State of Cloud"
report in a GitOps manner. It provides:
- &#9989; Cloud codification, identify un-managed resources and generate corresponding Terraform code and import statements/import blocks

- &#9989; Drift detection

- &#9989; Flag accounts creating changes outside your Terraform workflow

- &#9989; Whole-cloud cost estimation, powered by Infracost

- &#9989; Whole-cloud security scanning, powered by tfsec

## Quick Start
### All Cloud Provider Pre-requisites
1) Configure an environment variable file (use one of our [templates](https://github.com/dragondrop-cloud/cloud-concierge/tree/dev/examples/environments/) to get started) to control the specifics of cloud-concierge's coverage.
2) Run `docker pull dragondropcloud/cloud-concierge:latest` to pull the latest image.

### AWS Quickstart
I) Run `aws configure` on your CLI and ensure that credentials with read-only access to your cloud are configured. If referencing state files stored in an s3 bucket, the credentials specified should be able to read those state files as well.

II) Run the cloud-concierge container using the following command:
   `
   docker run --env-file ./my-env-file.env -v main:/main -v ~/.aws:/main/credentials/aws:ro -w /main  dragondropcloud/cloud-concierge:latest
   `

If running on Windows, the substitute `$HOME/.aws:` for `~/.aws:` in the above command.

III) Check the Pull Request that has been created by cloud-concierge ([example output](https://github.com/dragondrop-cloud/cloud-concierge-example/pull/3)).

### Azure & GCP Quickstart
See more [here](https://docs.cloudconcierge.io/quick-start#gcp).

## How does it work?
1) cloud-concierge creates a representation of your cloud infrastructure as Terraform. Only read-only access should be given to cloud-concierge.
2) This representation is compared against your state files to detect drift, and identify resources outside of Terraform control.
3) Static security scans and cost estimation is performed on the Terraform representation.
4) Results and code are summarized in a [Pull Request](https://docs.cloudconcierge.io/how-it-works/pull-request-output) within the repository of your choice.

### NLP Engine
Usage of cloud-concierge is not tracked. To maintain a [slimmer final docker image](https://medium.com/@hello_9187/ripping-out-python-and-reducing-our-docker-image-size-by-87-b61beda90ce4)
we host the NLP model that recommends new Terraform resources to match with state files in a Python-based Google Cloud Function.

This modelling code does not log any of the anonymized data sent to it,
and is stored within this repository [here](https://github.com/dragondrop-cloud/cloud-concierge/blob/1c31a98cec6d2c1189c9e4da35c616de100c04bb/nlpengine/main.py).

## Contributing
Contributions in any form are highly encouraged. Check out our [contributing guide](CONTRIBUTING.md) to get started.

## Resources
- [Example Output](https://github.com/dragondrop-cloud/cloud-concierge-example/pull/3)
- [Documentation](https://docs.cloudconcierge.io)
