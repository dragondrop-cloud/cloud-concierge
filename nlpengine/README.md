## NLP Engine

### Function
When resources outside of Terraform control are identified, a typically manual step is to determine
to which state file a resource belongs. This is a potentially time-consuming process, especially
when there are many resources to be imported into many state files.

cloud-concierge automates this process by using a natural language processing (NLP) engine to
determine the state file to which a resource belongs. This way, when codified resources are output within
a Pull Request, they are output to an intelligently selected directory which corresponds to an appropriate state file.

### Methodology
We approach this task as a multi-label text-categorization problem. We map the problem space into the following concepts:
- **Categories**: A category corresponds to all the resources found within a state file. Generally, each state file will
have a logical grouping of resources, such as "Dev Networking" or "Prod App 2", and we ultimately need to assign codified resources
to a particular state file.
- **Documents**: A document is a resource that has been identified as being outside of Terraform control. We want to classify the documents
  (codified cloud resources) and assign them to a category (state file).

Of course, we train an NLP model on text data, not on JSON or HCL. As a result, cloud-concierge must first convert the JSON/HCL of state files and
newly codified resources into anonymized, representative, text data. For example the following resource:
```hcl
resource "aws_s3_bucket" "my_bucket" {
  bucket = "my-bucket"
  acl    = "private"
  region = "us-east-1"
  
  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}
```

Would be converted to text along the lines of:
```text
Resource my_bucket of type aws_s3_bucket in region us-east-1 with tag Name and value My bucket, tag Environment with value Dev.
```

### Why a Separate Endpoint
While keeping this NLP engine bundled within cloud-concierge directly, it yields the following problems:
- Including Python and spaCy dependencies in the cloud-concierge Docker image really [bloated](https://dragondrop.cloud/2023/09/19/ripping-out-python-and-reducing-our-docker-image-size-by-87/) the image size.
This makes it much slower to download and run cloud-concierge.
- Different users local hardware constraints may make it difficult to run the NLP engine locally in a reasonable amount of time.

So, by hosting the service in a dedicated, Google Cloud Function endpoint we can:
- Keep the cloud-concierge Docker image small and fast to download
- Speed up the complete cloud-concierge process by offloading the NLP engine to dedicated compute.

### Security
On top of receiving only highly anonymized "document" and "category" data to train on, the NLP engine runs in ephemeral
Google Cloud Functions and does not log or save any data. It is stateless and ephemeral.
