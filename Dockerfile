FROM golang:1.25rc2-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app
VOLUME /config
RUN go build -o main .
CMD ["/app/main"]
