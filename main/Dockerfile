# This dockerfile performs a multi-stage build.
# Stage 1) Creates a reference to the tfswitch binary.
# Stage 2) Creates a reference to the infracost binary.
# Stage 3) Creates a reference to the gcloud image.
# Stage 4) References the terraformer image.
# Stage 5) Builds an executable binary out of the cloud-concierge go-code.
# Stage 6) Places binaries within the gcloud container from stage 2 and executes
###################################################################################################
# 1) Reference to tfswitch binary
###################################################################################################
FROM golang:1.18.2-alpine3.15 as tfswitch
RUN apk update && apk add --no-cache bash curl git make
RUN curl -L https://raw.githubusercontent.com/warrensbox/terraform-switcher/release/install.sh | bash

###################################################################################################
# 2) Reference to tfsec binary
###################################################################################################
FROM aquasec/tfsec:v1.28.1 as tfsec

###################################################################################################
# 3) Reference to infracost binary
###################################################################################################
FROM infracost/infracost:0.10.19 as infracost

###################################################################################################
# 4) Building the gcloud cli with config-connector
###################################################################################################
FROM --platform=linux/amd64 gcr.io/google.com/cloudsdktool/google-cloud-cli:436.0.0-slim as gcloud

###################################################################################################
# 5) Reference to terraformer binary
###################################################################################################
FROM dragondropcloud/terraformer:latest as terraformer

###################################################################################################
# 6) Building the cloud-concierge executable
###################################################################################################
FROM golang:1.18.2-alpine3.15 as cloud-concierge
RUN apk update && apk add --no-cache bash git make

# Beginning to build the cloud-concierge executable
WORKDIR $GOPATH/src

COPY go.mod go.sum ./
RUN go mod download

COPY .. .

# Compiling the cloud-concierge executable
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
     go build -ldflags='-w -s -extldflags "-static"' -a \
     -o /go/bin/cloud-concierge .

###################################################################################################
# 8) Creating the final light-weight container that contains only the executables from previous steps.
###################################################################################################
FROM gcloud

# Installing python dependencies
RUN pip3 install --upgrade pip
COPY internal/python_scripts/requirements.txt python_scripts/requirements.txt
RUN pip3 install -r python_scripts/requirements.txt

# Code that changes most frequently is copied into the container last.
COPY --from=tfswitch /usr/local/bin/tfswitch /usr/local/bin/
COPY --from=terraformer /go/bin/terraformer /usr/local/bin/
COPY --from=infracost /usr/bin/infracost /usr/local/bin/
COPY --from=tfsec /usr/bin/tfsec /usr/local/bin/
COPY --from=cloud-concierge /go/bin/cloud-concierge /go/bin/cloud-concierge
COPY internal/python_scripts python_scripts

ENTRYPOINT ["/go/bin/cloud-concierge"]
