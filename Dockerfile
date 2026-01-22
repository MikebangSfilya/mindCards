FROM golang:1.24-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /mindcards ./cmd/app/main.go

FROM alpine:latest


RUN apk add --no-cache ca-certificates

WORKDIR /root/


COPY --from=builder /mindcards .


EXPOSE 8080

CMD ["./mindcards"]