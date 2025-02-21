FROM golang:1.24.0-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app
VOLUME /config
RUN go build -o main .
CMD ["/app/main"]
