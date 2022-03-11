use `grumbler_db`;

-- pass: aaaaaaaa
insert into `users` (`id`, `created_at`, `name`, `password`, `profile`) values
('test1', '2022-02-01 10:04:21', 'tEsT', '$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2', 'プロフィール中の\n改行コードが\nHTMLに反映されているか\n\nテスト'),
('test2', '2022-02-02 10:04:21', '尾骶骨テスト', '$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2', ' つちよし𠮷test '),
('test3', '2022-02-03 10:04:21', 'aAあ亜', '$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2', 'aAあ亜0123456789プロファイル');

INSERT INTO `grumbles` (`pk`, `content`, `user_id`, `created_at`) VALUES
('01FWMM3VMDVEG35G3ER9SG5G94', 'a', 'test1', '2022-02-24 10:04:21'),
('01FX1T8DPVG182CA9XJH392M47', 'l', 'test1', '2022-03-01 13:01:50'),
('01FX4VZHEN47XV7E3J6PFPDEJM', 'ああああああああああああ\naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\nafgre\nrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr', 'test1', '2022-03-02 17:29:39'),
('01FX5AEM3549N0937KGBV34D9J', 'aaaaaaaaaaaaaaaaaaaaa\nrrrrrrrrrrrrrrrrr\ngbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb', 'test1', '2022-03-02 21:42:33'),
('01FX73CT6H9P77JZ9RWDZ9X2SC', 'gea a', 'test1', '2022-03-03 14:17:43'),
('01FXH4V8WRNHPNZHH4CAY7KE6M', 'bkosrkbno', 'test2', '2022-03-02 11:55:30'),
('01FXH4VWEW6XPHP4M5PA3MBHT5', 'bisrnbsir\n\nojojoijijiji\n\nhtorskhos', 'test2', '2022-03-07 11:55:50');

INSERT INTO `follows` (`pk`, `created_at`, `src_user_id`, `dst_user_id`) VALUES
(3, '2022-02-13 10:04:21', 'test1', 'test3'),
(4, '2022-02-12 10:04:21', 'test1', 'test2'),
(12, '2022-02-14 10:04:21', 'test2', 'test1');
