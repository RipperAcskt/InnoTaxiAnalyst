CREATE TABLE IF NOT EXISTS innotaxi.users(
    id UUID,
    user_id Int64,
    name String,
    phone_number String,
    email String,
    raiting Float64,
) ENGINE = MergeTree() 
PRIMARY KEY (id)
ORDER BY (id);