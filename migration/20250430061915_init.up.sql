CREATE TABLE IF NOT EXISTS clients (
                                       id           text    PRIMARY KEY,
                                       capacity     int     NOT NULL,
                                       rate_per_sec int     NOT NULL
);

INSERT INTO clients VALUES ('default', 100, 10);