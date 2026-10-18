# syntax=docker/dockerfile:1

# ---- build stage: full Go toolchain, never shipped ----
FROM golang:1.27-alpine AS build
WORKDIR /src

# Dependencies first so this layer is cached until go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Static binary (no libc) so it runs on a distroless base image.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w \
      -X github.com/BHAV0207/Requip-assignment/internal/version.Version=${VERSION} \
      -X github.com/BHAV0207/Requip-assignment/internal/version.Commit=${COMMIT} \
      -X github.com/BHAV0207/Requip-assignment/internal/version.BuildTime=${BUILD_TIME}" \
    -o /out/server ./cmd/server

# ---- runtime stage: no shell, no package manager, runs as UID 65532 ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/server /server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
