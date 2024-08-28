CREATE TABLE dummy
(
    id       BIGINT         NOT NULL PRIMARY KEY,
    product  VARCHAR(255)   NOT NULL,
    price    DECIMAL(10, 2) NOT NULL,
    qty      BIGINT         NOT NULL,
    null_data VARCHAR(255),
    date     VARCHAR(255)   NOT NULL
)