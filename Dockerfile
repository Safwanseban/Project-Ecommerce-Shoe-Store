FROM golang:alpine AS builder
#maintainer info
LABEL maintainer="Safwanseban <safwansappu0@gmail.com>"
#installing git
RUN apk update && apk add --no-cache git
#current working directory
COPY templates /.
WORKDIR /PROJECT-ECOMMERCE
# Copy go mod and sum files
COPY go.mod .
COPY go.sum .


# COPY /templates ./templates/


# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download
COPY . .
#building the go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/


# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file

COPY --from=builder /PROJECT-ECOMMERCE/main .
COPY --from=builder /PROJECT-ECOMMERCE/.env .       

COPY . .
# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]

