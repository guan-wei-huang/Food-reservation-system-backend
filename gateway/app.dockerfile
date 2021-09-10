FROM golang:1.17-alpine
WORKDIR /order
COPY go.mod ./
COPY go.sum ./
RUN go mode download

COPY *.go ./
COPY order order
COPY restaurant restaurant
COPY user user
RUN go build -o /go/bin/app ./order/cmd/main.go
EXPOSE 7999
CMD [ "app" ]