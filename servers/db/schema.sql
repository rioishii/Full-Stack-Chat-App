create table if not exists users (
    id int not null auto_increment primary key,
    email nvarchar(320) not null,
    pass_hash char(60) not null,
    user_name varchar(255) not null,
    first_name varchar(64) not null,
    last_name varchar(128) not null,
    photo_url varchar(255) not null,
    UNIQUE(id),
    UNIQUE(user_name)
);