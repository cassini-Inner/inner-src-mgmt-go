CREATE TABLE IF NOT EXISTS GLOBALSKILLS
(
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_by  INTEGER REFERENCES USERS (id),
    value       VARCHAR(50) NOT NULL,
    time_created timestamptz NOT NULL default current_timestamp
)