drop table if exists category;

create table category (
    id serial not null primary key,
    category_name varchar(34) not null
);

drop table if exists author;

create table author (
    id serial not null primary key,
    author_name varchar(32) not null
);

drop table if exists book;

create table book(
    id serial not null primary key,
    book_name varchar(48),
    texts text not null,
    author_id integer references author(id),
    category_id integer references category(id) 
);