package main

import "time"

type Product struct {
	Fid      int
	Name     string
	Price    float32
	Quantity int
}

type Order struct {
	Id        int
	Rid       int
	Uid       int
	Products  *[]Product
	CreatedAt time.Time
}
