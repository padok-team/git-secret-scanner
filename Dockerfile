# Build git-secret-scanner binary
FROM docker.io/library/golang:1.23.3@sha256:d56c3e08fe5b27729ee3834854ae8f7015af48fd651cd25d1e3bcf3c19830174 AS builder

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

# Retrieve gitleaks binary
FROM ghcr.io/gitleaks/gitleaks:v8.21.2@sha256:0e99e8821643ea5b235718642b93bb32486af9c8162c8b8731f7cbdc951a7f46 AS gitleaks

# ---

# Retrieve trufflehog binary
FROM docker.io/trufflesecurity/trufflehog:v3.82.13@sha256:9abf17c8902d58c05d82f910cf5dec05d100912482e8002d88918511fb44b6f6 AS trufflehog

# ---

# Build the final image
FROM docker.io/library/alpine:3.20.3@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d

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
COPY --from=gitleaks --chmod=555 /usr/bin/gitleaks /usr/local/bin/gitleaks
COPY --from=trufflehog --chmod=555 /usr/bin/trufflehog /usr/local/bin/trufflehog

# Copy the binary to the production image from the builder stage
COPY --from=builder --chmod=511 /workspace/bin/git-secret-scanner /usr/local/bin/git-secret-scanner

# Use an unprivileged user
USER 65532:65532

# Run git-secret-scanner on container startup
ENTRYPOINT ["/usr/local/bin/git-secret-scanner"]
