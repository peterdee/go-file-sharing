FROM golang:1.24-alpine AS builder
WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server main.go

FROM alpine AS release
WORKDIR /release
COPY --from=builder /build/server /release/server
ENV IS_DOCKER_IMAGE="true"
EXPOSE 9000
RUN chmod +x server
ENTRYPOINT ["./server"]
