# Build stage
FROM golang:1.21 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
COPY --from=builder /app /app
EXPOSE 3000

ENTRYPOINT ["/app"]
