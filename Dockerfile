FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /tms ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup
RUN adduser -S appuser -G appgroup
USER appuser

WORKDIR /home/appuser

COPY --from=builder /app/static ./static
COPY --from=builder /tms .

EXPOSE 8080

ENTRYPOINT ["./tms"]
