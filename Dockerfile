FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /taskflow ./cmd/api

FROM alpine:3.22

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /taskflow /app/taskflow

EXPOSE 8080

CMD ["/app/taskflow"]
