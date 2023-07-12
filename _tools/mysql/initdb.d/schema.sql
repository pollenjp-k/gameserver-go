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
