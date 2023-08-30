package main

import "github.com/nats-io/stan.go"

func pub() {

	sc, err := stan.Connect("my-cluster", "client-2")
	if err != nil {
		panic("Нет коннекта к Stan")
	}
	defer sc.Close()
	sc.Publish("foo", []byte(one))
	sc.Publish("foo", []byte(two))
	sc.Publish("foo", []byte("{invalid_model}"))
	sc.Publish("foo", []byte("invalid_model+another_invalid_model"))
}

var one = `{
  "order_uid": "one",
  "track_number": "WB",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "model1",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WB",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
`

var two = `{
"order_uid": "two",
"track_number": "test_track_number",
"entry": "WBIL",
"delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
},
"payment": {
    "transaction": "model2",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
},
"items": [
    {
    "chrt_id": 9934930,
    "track_number": "test_track_number",
    "price": 453,
    "rid": "ab4219087a764ae0btest",
    "name": "Mascaras",
    "sale": 30,
    "size": "0",
    "total_price": 317,
    "nm_id": 2389212,
    "brand": "Vivienne Sabo",
    "status": 202
    }
],
"locale": "en",
"internal_signature": "",
"customer_id": "test",
"delivery_service": "meest",
"shardkey": "9",
"sm_id": 99,
"date_created": "2021-11-26T06:22:19Z",
"oof_shard": "1"
}`
