# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# It starts from the golang base image.
FROM golang

# Add maintainer info.
LABEL maintainer="Ícaro Ribeiro <icaroribeiro@hotmail.com>"

# Set the current working directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source from the current directory to the working directory inside the container.
COPY . .

# Build the Go app.
RUN go build -o main .

# Command to run the application.
CMD ["./main"]