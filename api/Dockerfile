# syntax=docker/dockerfile:1

# ---- build ----
FROM golang:1.24-alpine AS build
WORKDIR /src

# The bracket glob makes go.sum optional: we have no third-party deps yet, but
# this keeps the layer working once we add some.
COPY go.mod go.su[m] ./
RUN go mod download

COPY . .

ARG VERSION=dev
# CGO off + -trimpath gives a static, reproducible binary for a scratch image.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/server ./cmd/server

# ---- run ----
# distroless/static has no shell and no package manager, and already carries
# the CA bundle the Resend HTTPS call needs.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
