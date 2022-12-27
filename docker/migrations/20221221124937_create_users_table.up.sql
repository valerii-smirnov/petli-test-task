create table users
(
    id            uuid primary key      default uuid_generate_v4(),
    email         varchar(100) not null unique,
    password_hash varchar(255) not null,
    registered_at timestamp    not null default now()
);