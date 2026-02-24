FROM golang:1.25 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o driver ./cmd/driver

FROM alpine:3.23.2

WORKDIR /server

COPY --from=builder /build/driver ./driver

CMD ["/server/driver"]