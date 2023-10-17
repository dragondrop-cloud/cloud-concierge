# This dockerfile performs a multi-stage build.
# Stage 1) Creates a reference to the tfswitch binary.
# Stage 2) Creates a reference to the tfsec binary.
# Stage 3) Creates a reference to the infracost binary.
# Stage 4) Creates a reference to the terraformer binary.
# Stage 5) Builds an executable binary out of the cloud-concierge go-code.
# Stage 6) Places binaries within a base Alpine container and executes
###################################################################################################
# 1) Reference to tfswitch binary
###################################################################################################
FROM golang:1.21-alpine3.18 as tfswitch
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
# 4) Reference to terraformer binary
###################################################################################################
FROM dragondropcloud/terraformer:latest as terraformer

###################################################################################################
# 5) Building the cloud-concierge executable
###################################################################################################
FROM golang:1.21-alpine3.18 as cloud-concierge
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
# 6) Creating the final light-weight container that contains only the executables from previous steps.
###################################################################################################
FROM alpine:3.18.3
RUN apk update && apk add --no-cache gcompat

# Code that changes most frequently is copied into the container last.
COPY --from=tfswitch /usr/local/bin/tfswitch /usr/local/bin/
COPY --from=terraformer /go/bin/terraformer /usr/local/bin/
COPY --from=infracost /usr/bin/infracost /usr/local/bin/
COPY --from=tfsec /usr/bin/tfsec /usr/local/bin/
COPY --from=cloud-concierge /go/bin/cloud-concierge /go/bin/cloud-concierge

ENTRYPOINT ["/go/bin/cloud-concierge"]
