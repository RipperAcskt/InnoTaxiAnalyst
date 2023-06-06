CREATE TABLE IF NOT EXISTS innotaxi.orders(
    id UUID,
    order_id String,
    user_id String,
    driver_id UUID,
    driver_name String,
    driver_phone String,
    driver_rating Float64,
    taxi_type String,
    from String,
	to String,
	date String,
	status String
) ENGINE = MergeTree() 
PRIMARY KEY (id)
ORDER BY (id);
