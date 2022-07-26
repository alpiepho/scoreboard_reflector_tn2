FROM golang:1.18.4-alpine3.16
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build -o main ./...
CMD ["/app/main"]
