package main

import (
	"fmt"
	"time"
)

// Model :
type Model struct {
	receiver_phone string
	order_status   map[string]Value
	product_type   string
	province_code  int
	start_time     time.Time
	end_time       time.Time
}

// Value :
type Value struct {
	number           int
	money_collection int
}

func NewModel(receiver_phone string, order_status map[string]Value, product_type string, province_code int, start_time time.Time, end_time time.Time) Model {
	return Model{receiver_phone, order_status, product_type, province_code, start_time, end_time}
}

func (model Model) Print() {
	fmt.Print(model.receiver_phone, "\t", model.province_code, "\t", model.product_type, "\t", model.order_status, "\t", model.start_time, "\t", model.end_time, "\n")
}
