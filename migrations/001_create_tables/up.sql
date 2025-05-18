create table study_entity(
    id serial primary key,
    api_id int not null,
    kind varchar(64) not null,
    name varchar(64) not null
);

create table chat(
    id bigint primary key,
    kind varchar(64) not null,
    name varchar(64) not null,
    username varchar(64),
    study_entity_id int references study_entity(id),
    is_banned boolean not null
);

create table admin(
    id serial primary key,
    chat_id bigint not null references chat(id),
    level int not null
);
