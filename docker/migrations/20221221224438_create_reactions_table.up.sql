CREATE TYPE reaction_action AS ENUM ('like', 'dislike');

CREATE TABLE reactions
(
    liker_id   uuid references dogs (id) on delete cascade,
    liked_id   uuid references dogs (id) on delete cascade,
    action     reaction_action not null,
    constraint unique_liker_dog unique (liker_id, liked_id),
    created_at timestamp       not null default now()
)