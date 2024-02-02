create database if not exists test;
use test;

drop table if exists user;
CREATE TABLE user (
                      id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT ,
                      user_id BIGINT NOT NULL UNIQUE,
                      name VARCHAR(64) NOT NULL UNIQUE,
                      password VARCHAR(256) NOT NULL,
                      email VARCHAR(64) NOT NULL,
                      gender TINYINT NOT NULL,
                      create_at TIMESTAMP NOT NULL,
                      update_at TIMESTAMP NOT NULL
)engine=innoDB default charset=utf8mb4 collate=utf8mb4_general_ci;


drop table if exists community;
CREATE TABLE community (
                           id BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
                           community_id BIGINT UNSIGNED NOT NULL UNIQUE,
                           community_name VARCHAR(64) NOT NULL UNIQUE,
                           community_intro VARCHAR(256) NOT NULL,
                           create_at TIMESTAMP NOT NULL,
                           update_at TIMESTAMP NOT NULL
)engine=innoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

-- 创建一个post表
drop table if exists post;
CREATE TABLE post (
                      id BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
                      post_id BIGINT UNSIGNED NOT NULL UNIQUE,
                      title VARCHAR(256) NOT NULL,
                      content TINYTEXT NOT NULL,
                      author_id BIGINT UNSIGNED NOT NULL,
                      community_id BIGINT UNSIGNED NOT NULL,
                      status TINYINT NOT NULL DEFAULT 0 COMMENT '',
                      create_at TIMESTAMP NOT NULL,
                      update_at TIMESTAMP NOT NULL,
                      INDEX idx_community_id (community_id),
                      FOREIGN KEY (community_id) REFERENCES community(community_id) ON DELETE CASCADE
)engine=innoDB default charset=utf8mb4 collate=utf8mb4_general_ci;