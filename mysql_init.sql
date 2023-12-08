create database if not exists test;
use test;

drop table if exists `user`;
create table `user` (
                        `id` bigint(20) not null auto_increment,
                        `user_id` bigint(20) not null,
                        `name` varchar(64) not null,
                        `password` varchar(256) not null,
                        `email` varchar(64) not null,
                        `gender` tinyint(4) not null,
                        `create_at` timestamp null default current_timestamp,   -- 默认使用当前时间戳
                        `update_at` timestamp null default current_timestamp
                            on update current_timestamp,    -- 每次更新的时候使用当前时间戳
                        primary key (`id`),
                        unique key `index_user_id` (`user_id`),
                        unique key `index_name` (`name`)
) engine=innoDB default charset=utf8mb4 collate=utf8mb4_general_ci;


-- 创建一个账户表
drop table if exists `community`;
create table `community` (
                             `id` int(11) unsigned not null auto_increment,
                             `community_id` bigint(20) unsigned not null,
                             `community_name` varchar(64) not null collate utf8mb4_general_ci,
                             `community_intro` varchar(256) not null collate utf8mb4_general_ci,
                             `create_at` timestamp not null default current_timestamp,
                             `update_at` timestamp not null default current_timestamp on update current_timestamp,
                             primary key (`id`),
                             unique key `index_community_id` (`community_id`),
                             unique key `index_community_name` (`community_name`)
) engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

-- 创建一个post表
drop table if exists `post`;
create table `post` (
                        `id` bigint(20) unsigned not null auto_increment,
                        `post_id` bigint(20) unsigned not null,
                        `title` varchar(256) not null,
                        `content` tinytext not null,
                        `author_id` bigint(20) unsigned not null,
                        `community_id` bigint(20) unsigned not null,
                        `status` tinyint(4) not null default 0,
                        `create_at` timestamp not null default current_timestamp,
                        `update_at` timestamp not null default current_timestamp on update current_timestamp,
                        primary key (`id`),
                        unique key `index_post_id` (`post_id`)
) engine =InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;