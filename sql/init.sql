--sudo docker run --name postgresql -e POSTGRES_USER=taktashova -e POSTGRES_PASSWORD=password123  -p 5432:5432 -d postgres
drop table if exists "users" CASCADE;
drop table if exists threats CASCADE;
drop table if exists monitoring_requests CASCADE;
drop table if exists monitoring_requests_threats CASCADE;

create table "users"
(
    user_id           SERIAL unique           not null
        constraint user_pk
            primary key,
    login             varchar(40)             not null UNIQUE,
    is_admin          boolean   default false,
    name              varchar(40)             not null,
    password          varchar(64)             not null,
    registration_date timestamp default now() not null
);
--таблица услуг
create table threats
(
    threat_id   SERIAL unique not null
        constraint threat_pk
            primary key,
    name        varchar(60)   not null,
    summary     TEXT          not null,
    description TEXT          not null,
    image       TEXT          not null,
    count       int                    default 0,
    is_deleted  boolean       not null default FALSE,
    price       int
);
--таблица заявок
create table monitoring_requests
(
    request_id     SERIAL unique           not null
        constraint monitoring_requests_pk
            primary key,
    creator_id     int
        constraint creator_id_fk
            references "users" (user_id),
    status         varchar(20)             not null,
    creation_date  timestamp default now() not null,
    formation_date timestamp,
    ending_date    timestamp,
    admin_id       int
);
-- таблица связи м:м
create table monitoring_requests_threats
(
    id         serial not null,
    request_id int
        constraint request_id_fk
            references monitoring_requests (request_id),
    threat_id  int
        constraint threats_id_fk
            references threats (threat_id),
    unique (request_id, threat_id)
);

SELECT *
FROM threats
WHERE threat_id = 6;