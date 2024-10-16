-- +migrate Up
create table if not exists assets
(
    id              varchar(36) primary key not null,
    user_id         varchar(36) not null,
    symbol          varchar(36) not null,
    total_balance           decimal(40,20) default 0,
    freeze_balance          decimal(40,20) default 0,
    avail_balance       decimal(40,20) default 0,
    created_at      timestamptz not null default now(), 
    updated_at      timestamptz not null default now()
);


create table if not exists assets_logs
(
    id                      varchar(36) primary key not null,
    user_id                 varchar(36) not null,
    symbol                  varchar(36) not null,
    before_balance          decimal(40,20) default 0,
    amount                  decimal(40,20) default 0,
    after_balance           decimal(40,20) default 0,
    trans_id                varchar(36) not null,
    change_type             varchar(36) not null,
    info                    varchar(255),
    created_at              timestamptz not null default now(), 
    updated_at      timestamptz not null default now()
);

create table if not exists assets_freezes
(
    id              varchar(36) primary key not null,
    user_id         varchar(36) not null,
    symbol          varchar(36) not null,
    amount          decimal(40,20) default 0,
    freeze_amount   decimal(40,20) default 0,
    status          smallint not null default 0,
    trans_id        varchar(36) not null,
    freeze_type     varchar(36) not null,
    info            varchar(255),
    created_at      timestamptz not null default now(), 
    updated_at      timestamptz not null default now()
);



create unique index if not exists user_id_symbol_unique_index on assets (user_id, symbol);





-- +migrate Down

drop index if exists user_id_symbol_unique_index;
drop table if exists assets;
drop table if exists assets_logs;
drop table if exists assets_freezes;