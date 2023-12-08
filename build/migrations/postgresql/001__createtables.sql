CREATE TABLE IF NOT EXISTS link
(
    id SERIAL PRIMARY KEY,
    link text UNIQUE NOT NULL,
    link_short text UNIQUE NOT NULL
);

---- create above / drop below ----

DROP TABLE link;
