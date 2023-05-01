# Use an official Golang runtime as a parent image
FROM golang:1.16-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Set the environment variables
ENV MOVIE_SERVICE="movierec-ms" \
    MOVIE_ENV="local" \
    MOVIE_WEBSERVER_PORT="50081" \
    MOVIE_WEBSERVER_ROUTE_PREFIX="/v1" \
    MOVIE_MONGO_CONNECTION_STRING="mongodb+srv://utkarsh:utkuser@movierec.js2kxf5.mongodb.net/?retryWrites=true&w=majority" \
    MOVIE_MONGO_DATABASE="movierecommendation-ms"

RUN go mod download
RUN go get go.mongodb.org/mongo-driver/x/mongo/driver/ocsp@v1.11.4


# Build the Go app
RUN go build -o server ./cmd/server/main.go

# Expose port 50081
EXPOSE 50081

# Run the server command by default when the container starts
ENTRYPOINT ["./server"]
