FROM golang:1.20-alpine

# Install necessary build tools and dependencies
RUN apk add --no-cache git protobuf-dev make bash

# Set working directory
WORKDIR /app

COPY go.mod go.sum ./

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go mod download

# Set environment variables
ENV PATH="/go/bin:${PATH}"
ENV PS1="\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]\$ "

# Set the default command to bash
CMD ["/bin/bash"]