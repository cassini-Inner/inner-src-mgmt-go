CREATE TABLE IF NOT EXISTS USERS
(
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email       VARCHAR(100) UNIQUE NOT NULL,
    name        VARCHAR(50)         NOT NULL,
    role        VARCHAR(50)         NOT NULL,
    department  VARCHAR(50)         NOT NULL,
    bio         VARCHAR(5000)       NOT NULL,
    photoURL    VARCHAR(500)         NOT NULL,
    contact     VARCHAR(100)         NOT NULL,
    time_created timestamptz        NOT NULL default current_timestamp,
    time_updated timestamptz        NOT NULL default current_timestamp,
    isDeleted   BOOLEAN             NOT NULL default false
);
