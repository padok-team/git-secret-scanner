# Build git-secret-scanner binary
FROM docker.io/library/golang:1.27-alpine@sha256:fa3d596898d9f126af7c01d14ee7597bab3f39cccf392ef36b69fe77cd681f18 AS builder

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
FROM docker.io/library/alpine:3.24@sha256:3f4d6b32e26703c6082839234f43769d5046b3ac7c910033911d1e71d8cd630a AS runtime

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
COPY --from=docker.io/trufflesecurity/trufflehog:3.99.8@sha256:674b2f187620a19df77f44cf9a7d77efcc581a5662accef2f2739cbb46ba4a3a \
    --chmod=555 /usr/bin/trufflehog /usr/local/bin/trufflehog

# Copy the binary to the production image from the builder stage
COPY --from=builder --chmod=511 /workspace/bin/git-secret-scanner /usr/local/bin/git-secret-scanner

# Use an unprivileged user
USER 65532:65532

# Run git-secret-scanner on container startup
ENTRYPOINT ["/usr/local/bin/git-secret-scanner"]
