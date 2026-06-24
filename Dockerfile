FROM golang:1.26 as Builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:3.24

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=Builder /app/main .
RUN chmod +x ./main

EXPOSE 8080

CMD ["./main"]