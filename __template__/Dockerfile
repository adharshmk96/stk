# Start from the official Golang image to build our application.
FROM golang:1.21 AS build

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to download all dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the application. 
# This produces a statically linked executable by disabling cgo which 
# is not needed in a scratch container.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stktemplate .

# Now, start from a new stage with a minimal base image for a smaller final image.
FROM alpine:latest

# Copy the statically linked executable from the build stage to the current stage.
COPY --from=build /app/stktemplate .

# Expose the port the application listens on.
EXPOSE 8080

# Command to run the executable.
CMD ["./stktemplate"]
