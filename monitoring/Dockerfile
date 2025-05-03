FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /monitoring ./monitoring/main.go

FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=builder /monitoring /monitoring

ENTRYPOINT ["/monitoring"]