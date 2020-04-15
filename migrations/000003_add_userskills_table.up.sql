CREATE TABLE IF NOT EXISTS USERSKILLS
(
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id     INTEGER REFERENCES USERS (id),
    skill_id    INTEGER REFERENCES GLOBALSKILLS (id),
    time_created timestamptz NOT NULL default current_timestamp,
    time_updated timestamptz NOT NULL default current_timestamp,
    isDeleted   BOOLEAN     NOT NULL default false
)