CREATE TABLE DISCUSSIONS
(
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    job_id      INTEGER REFERENCES JOBS (id),
    created_by  INTEGER REFERENCES USERS (id),
    content     VARCHAR(150000) NOT NULL,
    time_created timestamptz     NOT NULL default current_timestamp,
    time_updated timestamptz     NOT NULL default current_timestamp,
    isDeleted   BOOLEAN         NOT NULL default false
)