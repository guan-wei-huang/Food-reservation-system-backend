FROM golang:1.17-alpine
WORKDIR /go/src/github.com/guan-wei-huang/reserve-restaurant
COPY go.mod go.sum ./
RUN go mod download

COPY user user
RUN GO111MODULE=on go build -o app ./user/cmd/main.go
EXPOSE 8080
CMD [ "./app" ]