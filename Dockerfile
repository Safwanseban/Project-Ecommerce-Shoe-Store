FROM golang:alpine AS builder
WORKDIR /PROJECT-ECOMMERCE
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .