FROM golang:1.19-alpine
WORKDIR /go/src/github.com/guan-wei-huang/reserve-restaurant
COPY go.mod go.sum ./
COPY wait-for-it.sh ./
RUN chmod 777 wait-for-it.sh
RUN go mod download
RUN apk update && apk add bash

COPY restaurant restaurant
RUN GO111MODULE=on go build -o app ./restaurant/cmd/main.go
EXPOSE 8080
CMD [ "./wait-for-it.sh", "restaurant_db:5432", "--", "./app"]
