-- Dummy schema
DROP TABLE IF EXISTS person;
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

DROP TABLE IF EXISTS place;
CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
);