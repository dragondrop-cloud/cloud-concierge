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
   <img src="https://img.shields.io/badge/docker-v0.1.10-blue" />
</a>

<a href="https://hub.docker.com/r/dragondropcloud/cloud-concierge" alt="Total Downloads">
   <img src="https://img.shields.io/badge/downloads-7.2k-maroon" />
</a>

<a href="https://cloud-concierge.slack.com/join/shared_invite/zt-1xx3sqsb6-cekIXs2whccZvbU81Xn5qg#/shared-invite/email" alt="Slack">
<img src="https://img.shields.io/badge/slack-Join_Us-blueviolet" />
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

- &#9989; Whole-cloud security scanning, powered by tfsec (checkov integration coming soon)

## In action (managed instance)
https://github.com/dragondrop-cloud/cloud-concierge/assets/52042939/f81f567c-2c04-4aaf-ba75-963c49bcfab5

## Quick Start
### All Cloud Provider Pre-requisites
1) Obtain an API token at https://app.dragondrop.cloud. For open source executions, we only collect data on when a cloud-concierge starts up (see the Telemetry section below).
2) Add the [cloud-concierge GitHub App](https://github.com/apps/cloud-concierge) to the repository into which generated Pull Requests should be output.
3) Configure an environment variable file (use one of our [templates](https://github.com/dragondrop-cloud/cloud-concierge/tree/dev/examples/environments/) to get started) to control the specifics of cloud-concierge's coverage.
4) Run `docker pull dragondropcloud/cloud-concierge:latest` to pull the latest image.

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
2) This representation is compared against your state files to detect drift, and identify resources outside of Terraform control
3) Static security scans and cost estimation is performed on the Terraform representation
4) Results and code are summarized in a [Pull Request](https://docs.cloudconcierge.io/how-it-works/pull-request-output) within the repository of your choice

### Telemetry
For OSS usage, Cloud Concierge only logs data to the dragondrop API whenever a container execution is started. This method can be viewed [here](main/internal/implementations/dragon_drop/http_dragondrop_oss_methods.go).
 
Jobs managed by the [dragondrop platform](https://dragondrop.cloud) log statuses over the course of the job execution and anonymized data for cloud visualizations to the dragondrop API. These methods
can be viewed [here](https://github.com/dragondrop-cloud/cloud-concierge/blob/dev/main/internal/implementations/dragon_drop/http_dragondrop_managed_execution.go) and
[here](https://github.com/dragondrop-cloud/cloud-concierge/blob/dev/main/internal/implementations/dragon_drop/http_dragondrop_managed_visualization.go).

## Contributing
Contributions in any form are highly encouraged. Check out our [contributing guide](CONTRIBUTING.md) to get started.

## Using at Scale with dragondrop.cloud
The cloud-concierge container is easy to manage in a single configuration.
If you are looking to use cloud-concierge at scale, however, the [dragondrop.cloud](https://dragondrop.cloud/how-it-works) management platform allows you to:
- Manage multiple cloud-concierge configurations through a user interface
- Manage different cron jobs for executing each configuration at desired intervals
- Consolidate multiple cloud-concierge executions into anonymized visualizations of drift, uncodified resources, cloud costs, and security risks.
- Continue to self-host cloud-concierge instances within your cloud using [serverless infrastructure](https://registry.terraform.io/namespaces/dragondrop-cloud).

## Resources
- [Example Output](https://github.com/dragondrop-cloud/cloud-concierge-example/pull/3)
- [Documentation](https://docs.cloudconcierge.io)
- [Roadmap](https://github.com/dragondrop-cloud/cloud-concierge/wiki/Roadmap)
- [Blog](https://medium.com/@hello_9187)
- [Slack](https://cloud-concierge.slack.com/join/shared_invite/zt-1xx3sqsb6-cekIXs2whccZvbU81Xn5qg#/shared-invite/email)
- [Tool Walk Through (low stakes and purely educational)](https://calendly.com/dragondrop-cloud/cloud-concierge-walk-through)
- [Managed Offering](https://docs.dragondrop.cloud/)
