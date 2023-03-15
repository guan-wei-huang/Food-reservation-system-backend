package main

import (
	"log"
	"reserve_restaurant/order"
)

func main() {
	config := order.GetConfig()

	// config.DatabaseDsn = "postgres://apple:123456@localhost:8892/apple?sslmode=disable"
	// config.Port = 3000

	r, err := order.NewOrderRepository(config.DatabaseDsn)
	if err != nil {
		log.Fatal("new order repo error: ", err)
		return
	}
	defer r.Close()

	s := order.NewService(r)
	log.Printf("Listening on port %v...", config.Port)
	log.Fatal(order.ListenGRPC(s, config))
}
