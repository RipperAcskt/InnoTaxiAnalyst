CREATE TABLE IF NOT EXISTS innotaxi.drivers(
    id UUID,
    name text,
    phone_number text,
    email text,
    password BLOB,
    raiting float,
    taxi_type text,
    status text,
    PRIMARY KEY((id))
);