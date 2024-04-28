# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o ./fire-go -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
COPY --from=builder /app/fire-go .
EXPOSE 3000
ENTRYPOINT ["./fire-go"]
