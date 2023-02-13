CREATE TABLE IF NOT EXISTS innotaxi.drivers(
    name text,
    phone_number text,
    email text,
    password text,
    rating float,
    status text,
    PRIMARY KEY((phone_number))
);