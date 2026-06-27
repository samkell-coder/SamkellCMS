FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
# THIS LINE UPDATED:
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o portfolio .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/portfolio .
COPY data/ data/
COPY static/ static/
COPY admin/ admin/
EXPOSE 8080
CMD ["./portfolio"]