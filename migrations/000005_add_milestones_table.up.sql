CREATE TABLE IF NOT EXISTS MILESTONES
(
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    assigned_to INT REFERENCES Users (id),
    job_id      INTEGER REFERENCES JOBS (id),
    title       VARCHAR(1000)  NOT NULL,
    description VARCHAR(10000) NOT NULL,
    duration    VARCHAR(50)    NOT NULL,
    resolution  VARCHAR(200)   NOT NULL,
    status      VARCHAR(20)    NOT NULL,
    time_created timestamptz    NOT NULL default current_timestamp,
    time_updated timestamptz    NOT NULL default current_timestamp,
    isDeleted   BOOLEAN        NOT NULL default false
)