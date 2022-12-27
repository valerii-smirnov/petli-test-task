CREATE TYPE dog_sex AS ENUM ('male', 'female');

CREATE TABLE dogs
(
    id         uuid primary key       default uuid_generate_v4(),
    user_id    uuid references users (id) on delete cascade,
    name       varchar(30)   not null,
    sex        dog_sex       not null,
    age        integer       not null,
    breed      varchar(30)   not null,
    image      varchar(1024) not null,
    created_at timestamp     not null default now(),
    updated_at timestamp     not null default now()
)