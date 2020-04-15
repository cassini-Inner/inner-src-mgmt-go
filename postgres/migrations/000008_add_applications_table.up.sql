CREATE TABLE applications
(
    id           INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    milestone_id INTEGER REFERENCES MILESTONES (id),
    applicant_id INTEGER REFERENCES USERS (id),
    status       VARCHAR(20) NOT NULL default 'pending',
    notes        VARCHAR(1000),
    time_created  timestamptz NOT NULL default current_timestamp,
    time_updated  timestamptz NOT NULL default current_timestamp
)