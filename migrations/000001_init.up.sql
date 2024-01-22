CREATE SCHEMA IF NOT EXISTS the_name_service;

CREATE TABLE IF NOT EXISTS the_name_service.humans (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,
    patronymic TEXT NULL,
    age INT NOT NULL CHECK ( 0 <= age AND age <= 127 ),
    nation VARCHAR(2) NOT NULL,
    sex VARCHAR(6) NOT NULL,   
    created_at TIMESTAMP NOT NULL DEFAULT(now() AT TIME ZONE 'utc'),
    last_update TIMESTAMP NOT NULL DEFAULT(now() AT TIME ZONE 'utc')
);

ALTER TABLE the_name_service.humans
ADD CONSTRAINT FCS_uni UNIQUE NULLS NOT DISTINCT (name, surname, patronymic);

