MicroService: 
1. reserve service (預約服務)
4. user service (用戶服務)
5. shopper service (商店服務)

2. place service (位置服務)
3. payment service (支付服務)



---
go.mod && go.sum 的部分可以在 dockerfile 中做複製
(COPY ...)

---
compile .proto  file
`protoc --go_out=plugins=grpc:. *.proto`

---
proto style:
1. Use camel as func name (ex. NewUser())
2. Use _ as field name (ex. string favorite_song)