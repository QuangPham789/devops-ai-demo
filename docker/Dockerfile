# docker/Dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /src

# copy go.mod và go.sum từ thư mục app (không có ../)
COPY app/go.mod app/go.sum ./
RUN go mod download

# copy toàn bộ mã nguồn
COPY app/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:3.18
RUN adduser -D -g '' appuser
USER appuser
WORKDIR /app
COPY --from=builder /src/server .
EXPOSE 8080
CMD ["./server"]
