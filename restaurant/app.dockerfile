FROM golang:1.17-alpine
WORKDIR /go/src/github.com/guan-wei-huang/reserve-restaurant
COPY go.mod go.sum ./
COPY wait-for-it.sh ./
RUN chmod 777 wait-for-it.sh
RUN go mod download
RUN apk update && apk add bash

COPY restaurant restaurant
RUN GO111MODULE=on go build -o app ./restaurant/cmd/main.go
EXPOSE 8001
CMD [ "./wait-for-it.sh", "restaurant_db:5432", "--", "./app"]
