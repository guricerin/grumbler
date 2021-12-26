drop database if exists `grumbler_db`;
create database `grumbler_db`;

drop user if exists `grumbler`@`%`;
create user `grumbler`@`%` identified by 'grumbler1234';
grant all on `grumbler_db`.* to `grumbler`@`%`;
flush privileges;

use `grumbler_db`;

drop table if exists `users`;
drop table if exists `sessions`;
drop table if exists `grumbles`;

create table `users` (
    `pk` int unsigned auto_increment primary key not null,
    `id` varchar(255) unique not null,
    `name` varchar(255) not null,
    `password` varchar(255) not null,
    `profile` varchar(255) not null
);

create table `sessions` (
    `pk` int unsigned auto_increment primary key not null,
    `token` text not null,
    `user_pk` int unsigned not null
);

create table `grumbles` (
    `pk` varchar(26) primary key not null, -- ulid
    `content` varchar(300) not null,
    `user_id` varchar(255) not null,
    `created_at` datetime not null -- timestampは2038年問題があるからボツ
);