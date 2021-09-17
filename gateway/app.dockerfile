FROM golang:1.17-alpine
WORKDIR /go/src/github.com/guan-wei-huang/reserve-restaurant
COPY go.mod go.sum ./
RUN go mod download

COPY order order
COPY gateway gateway
COPY user user
COPY restaurant restaurant
RUN GO111MODULE=on go build -o app ./gateway
EXPOSE 8080
CMD [ "./app" ]