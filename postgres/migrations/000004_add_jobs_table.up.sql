CREATE TABLE IF NOT EXISTS JOBS
(
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    createdBy   INT REFERENCES USERS (id),
    title       VARCHAR(1000)  NOT NULL,
    description VARCHAR(10000) NOT NULL,
    difficulty  VARCHAR(20)    NOT NULL,
    status      VARCHAR(20)    NOT NULL default 'open',
    time_created timestamptz     NOT NULL default current_timestamp,
    time_updated timestamptz     NOT NULL default current_timestamp,
    isDeleted   BOOLEAN        NOT NULL default false
)
