FROM golang:1.22.7-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app
VOLUME /config
RUN go build -o main .
CMD ["/app/main"]
