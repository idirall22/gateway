# grpc server image
# stage 1
FROM golang:1.14 as builder

LABEL maintainer="Idir"

WORKDIR /app

# copy all files
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server.go

# stage 2
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server .

# Command to run the executable
CMD ["./server"]