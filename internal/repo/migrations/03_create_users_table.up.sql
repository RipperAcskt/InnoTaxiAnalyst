CREATE TABLE IF NOT EXISTS innotaxi.users(
    id UUID,
    user_id Int64,
    name String,
    phone_number String,
    email String,
    rating Float64,
    num_of_marks Int64
) ENGINE = MergeTree() 
PRIMARY KEY (id)
ORDER BY (id);