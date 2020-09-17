FROM golang:latest

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers

RUN go build -o main .

CMD ["/app/main"]

EXPOSE 8080
