create table if not exists notifications
(
    id           serial primary key,
    recipient_id int references users (id)             not null,
    sender_id    int references users (id)             not null,
    type         varchar(50)                           not null default '',
    read         bool        default false             not null,
    job_id       int references jobs (id)              not null,
    time_created timestamptz default CURRENT_TIMESTAMP not null
);