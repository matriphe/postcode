FROM golang:alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /postcode-smoke ./.github/ci/go

FROM busybox:latest

COPY --from=builder /postcode-smoke /postcode-smoke
ENTRYPOINT ["/postcode-smoke"]
