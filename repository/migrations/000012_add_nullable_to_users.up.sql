alter table users
alter column email drop not null,
alter column name drop not null,
alter column role drop not null,
alter column department drop not null,
alter column bio drop not null,
alter column photo_url drop not null,
alter column contact drop not null,
alter column time_created  drop not null,
alter column time_updated  drop not null,
alter column is_deleted  drop not null,
alter column github_url  drop not null,
alter column onboarded  drop not null;
