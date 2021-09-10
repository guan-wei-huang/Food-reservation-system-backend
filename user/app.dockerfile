FROM golang:1.17-alpine
WORKDIR /user
COPY go.mod ./
COPY go.sum ./
RUN go mode download

COPY *.go ./
RUN go build -o /go/bin/app ./user/cmd/main.go
EXPOSE 8002
CMD [ "app" ]