CREATE TABLE IF NOT EXISTS innotaxi.drivers(
    id UUID,
    driver_id UUID,
    name String,
    phone_number String,
    email String,
    raiting Float64,
    taxi_type String,
) ENGINE = MergeTree() 
PRIMARY KEY (id)
ORDER BY (id);