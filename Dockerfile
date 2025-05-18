FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/app"]