FROM golang:1.17-alpine
WORKDIR /go/src/github.com/guan-wei-huang/reserve-restaurant
COPY go.mod go.sum ./
RUN go mod download

COPY order order
RUN GO111MODULE=on go build -o app ./order/cmd/main.go
EXPOSE 8000
CMD [ "./app" ]