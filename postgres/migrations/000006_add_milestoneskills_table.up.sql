CREATE TABLE IF NOT EXISTS MILESTONESKILLS
(
    id           INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    milestone_id INTEGER REFERENCES MILESTONES (id),
    skill_id     INTEGER REFERENCES GLOBALSKILLS (id),
    time_created  timestamptz NOT NULL default current_timestamp,
    time_updated  timestamptz NOT NULL default current_timestamp,
    is_deleted BOOLEAN NOT NULL default false
)