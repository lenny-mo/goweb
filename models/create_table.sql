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
