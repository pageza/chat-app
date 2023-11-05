# Use the official Golang image to create a build artifact.
FROM golang:1.18 as builder

# Copy local code to the container image.
WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server ./main.go

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]
