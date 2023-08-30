package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func InsertData(conn *pgx.Conn, order OrderInfo) { //записываем данные из структуры в базу данных
	err := InsertDataPayment(conn, order)
	if err != nil {
		fmt.Println(err)
	}
	uid, err := InsertDataOrder(conn, order)
	if err != nil {
		fmt.Println(err)
	}
	id, err := InsertDataDelivery(conn, order)
	if err != nil {
		fmt.Println(err)
	}
	err = InsertOrderDelivery(conn, uid, id)
	if err != nil {
		fmt.Println(err)
	}
	err = InsertDataItems(conn, order)
	if err != nil {
		fmt.Println(err)
	}
}

func InsertDataOrder(conn *pgx.Conn, order OrderInfo) (order_uid string, err error) {
	query := `insert into order_info (order_uid, track_number, entry, locale, internal_signature,customer_id, delivery_service, shardkey, sm_id, date_created,oof_shard)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) returning order_uid`
	if err = conn.QueryRow(context.Background(), query, order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService,
		order.Shardkey, order.SmId, order.DateCreated, order.OofShard).Scan(&order_uid); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			//panic_msg:=fmt.Sprintf("Не удалось вставить данные :%s",pgErr.Message)
			//panic(panic_msg)
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
		}
	}
	return order_uid, nil
}

func InsertDataDelivery(conn *pgx.Conn, order OrderInfo) (id int, err error) {
	query := `insert into deliveries (name, phone, zip, city, address, region, email)
		values ($1,$2,$3,$4,$5,$6,$7) returning id`
	if err = conn.QueryRow(context.Background(), query, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region,
		order.Delivery.Email).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
		}
	}
	return id, nil
}

func InsertDataPayment(conn *pgx.Conn, order OrderInfo) (err error) {
	query := `insert into payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
		values ($1,$2,$3,$4,$5,$6,$7, $8, $9, $10)`
	if err = conn.QueryRow(context.Background(), query, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee).Scan(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
		}
	}
	return nil
}

func InsertOrderDelivery(conn *pgx.Conn, order_uid string, id int) (err error) { // связываем таблицы order_info и delivery в базе данных
	query := `insert into order_delivery (order_uid,delivery_id)
    	values ($1,$2)`
	if err = conn.QueryRow(context.Background(), query, order_uid, id).Scan(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
		}
	}
	return nil
}

func InsertDataItems(conn *pgx.Conn, order OrderInfo) (err error) {
	for i := 0; i < len(order.Items); i++ {
		query := `insert into items (chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status)
			values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
		if err = conn.QueryRow(context.Background(), query, order.Items[i].ChrtId, order.Items[i].TrackNumber, order.Items[i].Price, order.Items[i].Rid, order.Items[i].Name, order.Items[i].Sale, order.Items[i].Size,
			order.Items[i].TotalPrice, order.Items[i].NmId, order.Items[i].Brand, order.Items[i].Status).Scan(); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				fmt.Println(pgErr.Message) // => syntax error at end of input
				fmt.Println(pgErr.Code)    // => 42601
			}
		}
	}
	return nil
}

func GetDataByUid(conn *pgx.Conn, order_uid string) (string, error) { // Получаем данные по order_uid из базы данных
	var order OrderInfo
	query := `select oi.*, to_jsonb(p.*) as "payment", (select jsonb_agg((to_jsonb(i.*))) from items i )  as "items", to_jsonb(del) as "delivery" 
		from order_info oi left join payments p on p."transaction" = oi.order_uid left join items i on i.track_number = oi.track_number
		join (select d.name,d.phone ,d.zip ,d.city ,d.address ,d.region ,d.email from deliveries d where d.id =(select od.delivery_id from order_delivery od where od.order_uid=$1)
		) as del on true where oi.order_uid = $1 limit 1`
	if err := conn.QueryRow(context.Background(), query, order_uid).Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId,
		&order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard, &order.Payment, &order.Items, &order.Delivery); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
		}
	}
	res, _ := json.Marshal(&order)
	return string(res), nil
}

func GetOrderUid(conn *pgx.Conn) (slice_uid []string) {
	query := `select array_agg(order_uid) from order_info`
	err := conn.QueryRow(context.Background(), query).Scan(&slice_uid)
	if err != nil {
		fmt.Println(err)
	}
	return slice_uid
}

func InsertInvalidData(conn *pgx.Conn, data string) (err error) {
	query := `insert into invalid_data(data) values ($1)`
	if err = conn.QueryRow(context.Background(), query, data).Scan(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
		}
	}
	return nil
}

var sqlconfig = `drop table if exists items,order_info, payments,deliveries, order_delivery, invalid_data;
	create table deliveries
	(
		id      serial primary key,
		name    varchar(64),
		phone   varchar(64),
		zip     varchar(64),
		city    varchar(64),
		address varchar(64),
		region  varchar(64),
		email   varchar(64)
	);
	create table payments
	(
		transaction   varchar(64) primary key,
		request_id    varchar(64),
		currency      varchar(64),
		provider      varchar(64),
		amount        integer,
		payment_dt    bigint,
		bank          varchar(64),
		delivery_cost integer,
		goods_total   integer,
		custom_fee    integer
	);
	create table order_info
	(
		order_uid          varchar(64) primary key references payments (transaction),
		track_number       varchar(64) unique,
		entry              varchar(64),
		locale             varchar(64),
		internal_signature varchar(64),
		customer_id        varchar(64),
		delivery_service   varchar(64),
		shardkey           varchar(64),
		sm_id              integer,
		date_created       timestamp,
		oof_shard          varchar(64)
	);
	create table items
	(
		chrt_id      integer,
		track_number varchar(64) references order_info (track_number),
		price        integer,
		rid          varchar(64),
		name         varchar(64),
		sale         integer,
		size         varchar(64),
		total_price  integer,
		nm_id        integer,
		brand        varchar(64),
		status       integer
	);
	create table order_delivery
	(
		order_uid   varchar(64) references order_info (order_uid),
		delivery_id int references deliveries (id)
	);
	create table invalid_data
	(
		id        serial primary key,
		data      varchar,
		timestamp timestamp default now()
	);`
