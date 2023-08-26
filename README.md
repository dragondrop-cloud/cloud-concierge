# cloud-concierge
<p align="center">
<img width="500" src=./images/cloud-concierge-logo.png>
</p>
<h2 align="center">
<a href="https://github.com/dragondrop-cloud/cloud-concierge-example/pull/2" target="_blank">Example Output</a> |
<a href="https://docs.cloudconcierge.io" target="_blank">Docs</a>
</h2>

## Why cloud-concierge?
cloud-concierge is a container that integrates with your existing Terraform management stack.
All results and codified resources are output via a digestible [Pull Request](https://github.com/dragondrop-cloud/cloud-concierge-example/pull/2) to a repository of your choice, providing you with a "State of Cloud"
report in a GitOps manner. It provides:
- &#9989; Cloud codification, identify un-managed resources and generate corresponding Terraform code and import statements/import blocks

- &#9989; Drift detection

- &#9989; Flag accounts creating changes outside your Terraform workflow

- &#9989; Whole-cloud cost estimation, powered by Infracost

- &#9989; Whole-cloud security scanning, powered by tfsec (checkov integration coming soon)

## In action (managed instance)
https://github.com/dragondrop-cloud/cloud-concierge/assets/52042939/f81f567c-2c04-4aaf-ba75-963c49bcfab5

## Motivation
Many teams build their own Terraform management "stacks" using major cloud provider state backends
and tools like Atlantis for running `plan` and `apply` and state-locking. 

For more sophisticated tooling, some may turn to tools like Terraform Cloud,
Scalr, Spacelift and Firefly. We find, however, that these tool's pricing can become particularly onerous
when wanting to self-host runners or access the most desired features like drift detection, security scanning, etc.

## Quick Start
### All Cloud Provider Pre-requisites
0) Obtain an API token at https://app.dragondrop.cloud. We only collect data on when a cloud-concierge starts up (this can be verified here).
1) Configure an environment variable file (use one of our templates to get started) to control the specifics of cloud-concierge's coverage.
2) Make sure you have Docker available on your local machine.

### AWS Quickstart
I) Run aws config on your CLI and ensure that credentials with read-only access to your cloud are configured. If referencing state files stored in an s3 bucket, the credentials specified should be able to read those state files as well.

II) Run the cloud-concierge container using the following command:
   `
   docker run --env-file ./my-env-file.env -v main:/main -v ~/.aws:/main/credentials/aws:ro -w /main  dragondropcloud/cloud-concierge:latest
   `

III) Upon job completion, check the repository against which you configured cloud-concierge to run. There will be a new Pull Request that has been created by Cloud Concierge. [Example Output](https://github.com/dragondrop-cloud/cloud-concierge-example/pull/2).

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
- [Example Output](https://github.com/dragondrop-cloud/cloud-concierge-example/pull/2)
- [Documentation](https://docs.cloudconcierge.io)
- [Roadmap](https://github.com/dragondrop-cloud/cloud-concierge/wiki/Roadmap)
- [Blog](https://medium.com/@hello_9187)
- [Slack](https://cloud-concierge.slack.com/join/shared_invite/zt-1xx3sqsb6-cekIXs2whccZvbU81Xn5qg#/shared-invite/email)
- [Tool Walk Through (low stakes and purely educational)](https://calendly.com/dragondrop-cloud/cloud-concierge-walk-through)
- [Managed Offering](https://docs.dragondrop.cloud/)
