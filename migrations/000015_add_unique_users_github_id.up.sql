alter table users
add constraint unique_github_id unique (github_id);