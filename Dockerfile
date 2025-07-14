# Copyright (c) 2024-2025 Joshua Sing <joshua@joshuasing.dev>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# Build stage
FROM golang:1.24.5-alpine3.22@sha256:ddf52008bce1be455fe2b22d780b6693259aaf97b16383b6372f4b22dd33ad66 AS builder

# Add ca-certificates, timezone data
RUN apk --no-cache add --update ca-certificates tzdata

# Create non-root user
RUN addgroup --gid 65532 starlink_exporter && \
    adduser  --disabled-password --gecos "" \
    --home "/etc/starlink_exporter" --shell="/sbin/nologin" \
    -G starlink_exporter --uid 65532 starlink_exporter

WORKDIR /build/starlink_exporter
COPY . .

RUN go mod download
RUN go mod verify

RUN GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) CGO_ENABLED=0 GOGC=off \
    go build -trimpath -ldflags "-s -w" -o /build/starlink_exporter/dist/starlink_exporter ./cmd/starlink_exporter

## Run stage
FROM scratch

# Build metadata
ARG VERSION
ARG VCS_REF
ARG BUILD_DATE

LABEL maintainer="Joshua Sing <joshua@joshuasing.dev>"
LABEL org.opencontainers.image.created=$BUILD_DATE \
      org.opencontainers.image.authors="Joshua Sing <joshua@joshuasing.dev>" \
      org.opencontainers.image.url="https://github.com/joshuasing/starlink_exporter" \
      org.opencontainers.image.source="https://github.com/joshuasing/starlink_exporter" \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.vendor="Joshua Sing <joshua@joshuasing.dev>" \
      org.opencontainers.image.title="Starlink Prometheus Exporter" \
      org.opencontainers.image.description="A simple Starlink exporter for Prometheus"

# Copy files
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/starlink_exporter /etc/starlink_exporter
COPY --from=builder /build/starlink_exporter/dist/starlink_exporter /usr/local/bin/starlink_exporter

USER starlink_exporter:starlink_exporter
EXPOSE 9451/tcp
ENTRYPOINT ["/usr/local/bin/starlink_exporter"]
