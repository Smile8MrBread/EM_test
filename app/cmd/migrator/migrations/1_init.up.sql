CREATE TABLE IF NOT EXISTS Songs
(
    id      serial primary key ,
    song    character varying(256) NOT NULL ,
    squad character varying(256) NOT NULL ,
    text    character varying(256)
)