create table if not exists reviews
(
    id           serial      not null unique primary key,
    rating       int         not null,
    remark       varchar(1000)        default '' not null,
    milestone_id int references milestones (id),
    user_id      int         not null references users (id),
    time_created timestamptz not null default current_timestamp,
    time_updated timestamptz not null default current_timestamp,
    is_deleted   boolean default false
);