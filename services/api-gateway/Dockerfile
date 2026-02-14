FROM golang:1.25-bookworm AS build-base

WORKDIR /app

COPY go.sum go.mod ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .

RUN go build \
    -ldflags="-linkmode external -extldflags -static" \
    -tags netgo \
    -o api-gateway ./cmd


# Use separate stage for deployable image
FROM scratch

COPY --from=build-base /app/api-gateway api-gateway

EXPOSE 8080

CMD ["/api-gateway"]
