FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /scheduler ./cmd/scheduler

FROM alpine:latest
COPY --from=builder /scheduler /scheduler
EXPOSE 9090
CMD ["/scheduler"]