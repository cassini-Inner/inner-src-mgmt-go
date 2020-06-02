alter table users
    alter column email set not null,
    alter column name set not null,
    alter column role set not null,
    alter column department set not null,
    alter column bio set not null,
    alter column photo_url set not null,
    alter column contact set not null,
    alter column time_created set not null,
    alter column time_updated set not null,
    alter column is_deleted set not null,
    alter column github_url set not null,
    alter column onboarded set not null;
