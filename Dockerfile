FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o sentinel ./cmd/sentinel

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/sentinel .
COPY *.yaml ./

EXPOSE 8080
ENTRYPOINT ["./sentinel"]
