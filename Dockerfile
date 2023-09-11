FROM golang:1.21-alpine AS builder

WORKDIR /go/src/app

COPY ./src .

# Build the Go app
RUN go build -o color-generator .

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/app/color-generator /app/
COPY ./app/static /app/static

RUN mkdir /app/data
RUN chmod +x ./color-generator

EXPOSE 8080

# Command to run the executable when the container starts
CMD ["./color-generator"]
