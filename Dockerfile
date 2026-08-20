# Этап 1: Сборка
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
COPY pkg/ ./pkg/

COPY *.db ./

COPY web/ ./web/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/todo .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo .

COPY --from=builder /app/*.db ./

COPY --from=builder /app/web ./web

EXPOSE 7540
CMD [ "./todo" ]