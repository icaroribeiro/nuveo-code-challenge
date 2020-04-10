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
RUN go build -o main

# I could verify that it takes sometime for the database container to get up.
# This way, we can't connect to the database at the start of the program.

# To solve this problem, we could apply one of the following approaches:

# i) Set the container to restart on-failure in docker-compose file.
CMD ["./main"]

# ii) Use a bash script to wait for the database to get up before running the program.
# For that, use wait-for-it.sh to test and wait on the availability of a TCP host and port.

# RUN git clone https://github.com/vishnubob/wait-for-it.git

# CMD ./wait-for-it/wait-for-it.sh --host=${DB_HOST} --port=${DB_PORT} --timeout=60 -- /app/main