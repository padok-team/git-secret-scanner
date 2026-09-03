# Build git-secret-scanner binary
FROM docker.io/library/golang:1.27-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder

ARG TARGETOS
ARG TARGETARCH
ARG PACKAGE=github.com/padok-team/git-secret-scanner
ARG VERSION
ARG COMMIT_HASH
ARG BUILD_TIMESTAMP

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY internal/ internal/
COPY main.go main.go

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a \
    -ldflags="\
    -X ${PACKAGE}/cmd.Version=${VERSION} \
    -X ${PACKAGE}/cmd.CommitHash=${COMMIT_HASH} \
    -X ${PACKAGE}/cmd.BuildTimestamp=${BUILD_TIMESTAMP}" \
    -o bin/git-secret-scanner main.go

# ---

# Build the final image
FROM docker.io/library/alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b AS runtime

WORKDIR /home/git-secret-scanner

ENV UID=65532
ENV GID=65532
ENV USER=git-secret-scanner
ENV GROUP=git-secret-scanner

# Install required packages
RUN apk update --no-cache
RUN apk add --no-cache \
    bash \
    git \
    binutils

# Create a non-root user to run the app
RUN addgroup -g $GID $GROUP
RUN adduser \
    --disabled-password \
    --no-create-home \
    --home $(pwd) \
    --uid $UID \
    --ingroup $GROUP \
    $USER

# Copy the scanners to the production image from the scanners stage
COPY --from=ghcr.io/gitleaks/gitleaks:v8.30.1@sha256:c00b6bd0aeb3071cbcb79009cb16a60dd9e0a7c60e2be9ab65d25e6bc8abbb7f \
    --chmod=555 /usr/bin/gitleaks /usr/local/bin/gitleaks
COPY --from=docker.io/trufflesecurity/trufflehog:3.97.4@sha256:562bc231afa9de3d04de44cfe624252b08207de1fc3cebc5e7ed92bed7f279e4 \
    --chmod=555 /usr/bin/trufflehog /usr/local/bin/trufflehog

# Copy the binary to the production image from the builder stage
COPY --from=builder --chmod=511 /workspace/bin/git-secret-scanner /usr/local/bin/git-secret-scanner

# Use an unprivileged user
USER 65532:65532

# Run git-secret-scanner on container startup
ENTRYPOINT ["/usr/local/bin/git-secret-scanner"]
