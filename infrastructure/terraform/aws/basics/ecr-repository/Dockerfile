# syntax=docker/dockerfile:1

ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.22

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

WORKDIR /app

COPY ./cmd/app/go.mod ./cmd/app/go.sum ./
RUN go mod download

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=bind,source=cmd/app,target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w" \
    -o /out/main \
    .

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /out/main /app/main

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["./main"]
