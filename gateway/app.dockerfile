FROM golang:1.17-alpine
WORKDIR /go/src/github.com/guan-wei-huang/reserve-restaurant
COPY go.mod go.sum ./
COPY wait-for-it.sh ./
RUN go mod download
RUN chmod 777 wait-for-it.sh
RUN apk update && apk add bash

COPY order order
COPY gateway gateway
COPY user user
COPY restaurant restaurant
RUN GO111MODULE=on go build -o app ./gateway
EXPOSE 8080
CMD [ "./wait-for-it.sh", "restaurant:3001", "--", "./app" ]