# syntax=docker/dockerfile:1

FROM golang:1.18

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./
COPY v1/ ./v1/
RUN ls -la

# Build
RUN CGO_ENABLED=0 GOARCH=arm GOOS=linux go build -o /iotgate-server
# RUN CGO_ENABLED=0 GOOS=linux go build -o /iotgate-server

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 80

# Run
CMD ["/iotgate-server"]