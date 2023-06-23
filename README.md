# cloud-concierge
<p align="center">
<img width="500" src=./images/cloud-concierge-logo.png>
</p>
<h2 align="center">
<a href="https://docs.cloudconcierge.io" target="_blank">Docs</a> |
<a href="https://https://www.youtube.com/watch?v=y8OSfQQMEL0&t=12s" target="_blank"> Managed Demo </a>
</h2>

Many teams build their own Terraform management "stacks" using major cloud provider state backends
and tools like Atlantis for running `plan` and `apply` and state-locking. 

For more sophisticated tooling, some may turn to tools like Terraform Cloud,
Scalr, Spacelift and Firefly. We find, however, that these tool's pricing can become particularly onerous
when wanting to self-host runners or access the most desired features like drift detection, security scanning, etc.

## Why Cloud Concierge?
Cloud Concierge is an open-sourced container that integrates with your existing Terraform management stack to provide:
- Cloud codification
- Drift Detection
- Flag accounts creating changes outside your Terraform workflow
- Whole-cloud-level cost estimation, powered by Infracost
- Whole-cloud-level security scanning, powered by tfsec (checkov integration coming soon)

All results and codified resources are output in a digestible Pull Request to a repository of your choice, providing you with a "State of Cloud"
report in a GitOps manner.

## How does it work?
1) Cloud Concierge creates a representation of your cloud infrastructure as Terraform
2) This representation is compared against your state files to detect drift, and identify resources outside of Terraform control
3) Static security scans and cost estimation is performed on the Terraform representation
4) Results and code are summarized in a Pull Request within the repository of your choice

## Getting Started
0) Retrieve an organization token from the dragondrop.cloud management platform [here](https://app.dragondrop.cloud).
1) Configure your environment variable file. This determines the execution behavior of the container. We provide example env configuration files for:
   - [AWS]()
   - [GCP]()
   - [Azure]()

Documentation on environment variables needed can be found [here]().

While Cloud Concierge validates environment variable formats upon start-up, we provide a UI for client-side validation of env vars
within the [dragondrop.cloud platform](https://app.dragondrop.cloud/env-var-validator) should faster iteration be desired.

4) Run the container with the following command:
```bash
docker run --env-file <path to env file> dragondrop-cloud/cloud-concierge:latest  # TODO: Volume specification needed
```

5) If using Terraform >= 1.5, Cloud Concierge generates [import blocks](https://medium.com/@hello_9187/terraform-1-5-xs-new-import-block-b8607c51287f) for newly codified resources directly.
If using Terraform < 1.5, we generate a `terraform import` command for each resource. These commands can be run manually,
or programmatically in a `plan` and `apply` manner using our [GitHub Action](https://github.com/dragondrop-cloud/github-action-tfstate-migration). 

### Running on a schedule
A common use case is to want to regularly scan for drift and un-codified resources. Cloud Concierge can easily be run
on a cron schedule using GitHub Actions. See our [example workflow]().

### Telemetry
For OSS usage, Cloud Concierge only logs whenever a container execution is started. This method can be viewed [here](pkg/implementations/dragon_drop/http_dragondrop_oss_methods.go).
 
Jobs managed by the [dragondrop platform](https://dragondrop.cloud) log statuses over the course of the job execution and anonymized data for cloud visualizations. These methods
can be viewed [here](pkg/implementations/dragon_drop/http_dragondrop_managed_execution.go) and [here]([here](pkg/implementations/dragon_drop/http_dragondrop_managed_execution.go)).

## Contributing
Contributions in any form are highly encouraged. Check out our [contributing guide](CONTRIBUTING.md) to get started.

## Other Resources
- [Documentation](https://docs.cloudconcierge.io)
- [Free Code Walk Through (low stakes and no pressure!)]()
- [Slack]()
- [Terraform Learning Resources](https://dragondrop.cloud/learn/terraform/)
- [Medium Blog](https://medium.com/@hello_9187)
- [Managed Offering](https://dragondrop.cloud/how-it-works/)
