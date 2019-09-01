FROM golang:1.12

WORKDIR /
COPY . .

RUN go get github.com/gorilla/mux
RUN go get go.mongodb.org/mongo-driver/bson
RUN go get go.mongodb.org/mongo-driver/mongo

EXPOSE 8080

CMD ["go","run","main.go"]
