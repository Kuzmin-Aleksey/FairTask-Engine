CREATE TYPE order_status AS ENUM ('processed', 'await', 'accept', 'reject');

CREATE TABLE orders
(
    id          INT PRIMARY KEY,
    parent_id   INT,
    text        TEXT                    NOT NULL,
    status      order_status            NOT NULL,
    executor_id INT,
    ts          TIMESTAMP DEFAULT now() NOT NULL
);


CREATE TABLE order_parameters
(
    order_id     INT  NOT NULL,
    parameter_id INT  NOT NULL,
    value        TEXT NOT NULL,
    UNIQUE (order_id, parameter_id)
);


CREATE TYPE value_type AS ENUM ('int', 'float', 'datetime', 'text', 'bool');

CREATE TABLE parameters
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type value_type   NOT NULL
);

CREATE TYPE executor_status AS ENUM ('active', 'inactive');

CREATE TABLE executors
(
    id     SERIAL PRIMARY KEY,
    name   VARCHAR(255)    NOT NULL,
    status executor_status NOT NULL
);

CREATE TABLE executor_parameters
(
    parameter_id INT  NOT NULL,
    executor_id  INT  NOT NULL,
    mask         TEXT NOT NULL,
    UNIQUE (parameter_id, executor_id)
)
