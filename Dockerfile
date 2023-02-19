FROM golang:1.17.2-alpine AS build

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the application
RUN go build -o /app/schoolapp

# Create a minimal Docker image for running the application
FROM alpine:3.14.2
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul
COPY --from=build /app/schoolapp /usr/local/bin/schoolapp
CMD ["schoolapp"]
