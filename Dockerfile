FROM golang:1.16-alpine

ENV HTTP_PORT=5000

RUN apk add --no-cache git

WORKDIR /go/src/app

COPY ./ ./

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o /main ./

# Expose port 5000 to the outside world
EXPOSE $HTTP_PORT

ENTRYPOINT ["/main"]

