--sudo docker run --name postgresql -e POSTGRES_USER=taktashova -e POSTGRES_PASSWORD=password123  -p 5432:5432 -d monitoring_app;
drop table if exists "users" CASCADE;
drop table if exists threats CASCADE;
drop table if exists monitoring_requests CASCADE;

create table "users"
(
    user_id           SERIAL                  not null
        constraint user_pk
            primary key,
    login             varchar(40)             not null,
    is_admin          boolean   default false,
    name              varchar(40)             not null,
    password_hash     varchar(64)             not null,
    registration_date timestamp default now() not null
);
--таблица услуг
create table threats
(
    threat_id   SERIAL      not null
        constraint threat_pk
            primary key,
    name        varchar(60) not null,
    description TEXT        not null,
    image       varchar(60) not null,
    count       int default 0,
    is_deleted  boolean     not null,
    price       int
);
--таблица заявок
create table monitoring_requests
(
    request_id     SERIAL                  not null,
    status         varchar(20)             not null,
    creation_date  timestamp default now() not null,
    formation_date timestamp,
    ending_date    timestamp,
    admin_id       int
        constraint monitoring_request_user_id_fk
            references "users" (user_id)
);