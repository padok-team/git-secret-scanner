# Build git-secret-scanner binary
FROM docker.io/library/golang:1.24.5@sha256:ef5b4be1f94b36c90385abd9b6b4f201723ae28e71acacb76d00687333c17282 AS builder

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
FROM ghcr.io/gitleaks/gitleaks:v8.28.0@sha256:cdbb7c955abce02001a9f6c9f602fb195b7fadc1e812065883f695d1eeaba854 AS gitleaks

# ---

# Retrieve trufflehog binary
FROM docker.io/trufflesecurity/trufflehog:3.90.3@sha256:f9a92af4d46ca171bffa5c00509414a19d9887c9ed4fe98d1f43757b52600e39 AS trufflehog

# ---

# Build the final image
FROM docker.io/library/alpine:3.22.1@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

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
