FROM golang:1.17-alpine
WORKDIR /restaurant
COPY go.mod ./
COPY go.sum ./
RUN go mode download

COPY *.go ./
RUN go build -o /go/bin/app ./restaurant/cmd/main.go
EXPOSE 8001
CMD [ "app" ]