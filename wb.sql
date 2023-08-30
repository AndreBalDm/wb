drop table if exists items,order_info, payments,deliveries, order_delivery, invalid_data;

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
);