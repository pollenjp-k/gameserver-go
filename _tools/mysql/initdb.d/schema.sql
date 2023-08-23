CREATE TABLE `user` (
  -- MySQL の AUTO_INCREMENT は 1 スタート
  -- 0 はアプリケーション側のint型のゼロ値であり、バリデーションで弾く実装を行う
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `token` varchar(255) DEFAULT NULL,
  `leader_card_id` int DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`)
) Engine=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';

CREATE TABLE `room` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  -- 楽曲ID
  `live_id` bigint NOT NULL,
  `host_user_id` bigint NOT NULL,
  `status` int NOT NULL DEFAULT 1,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `room_user` (
  `room_id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `live_difficulty` int NOT NULL,
  `status` int NOT NULL DEFAULT 1,
  PRIMARY KEY (`room_id`, `user_id`)
);


-- TODO: `room_user` と primary key の取り方を同じにしているが、他に良い方法はあるだろうか？
CREATE TABLE `score` (
  `room_id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `score` int NOT NULL,
  `judge_perfect` int NOT NULL,
  `judge_great` int NOT NULL,
  `judge_good` int NOT NULL,
  `judge_bad` int NOT NULL,
  `judge_miss` int NOT NULL,
  PRIMARY KEY (`room_id`, `user_id`)
);
