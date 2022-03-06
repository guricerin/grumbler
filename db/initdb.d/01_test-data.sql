use `grumbler_db`;

-- pass: aaaaaaaa
insert into `users` (`id`, `name`, `password`, `profile`) values ("test1", "tEsT", "$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2", "test profile");
insert into `users` (`id`, `name`, `password`, `profile`) values ("test2", "尾骶骨テスト", "$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2", " つちよし𠮷test ");
insert into `users` (`id`, `name`, `password`, `profile`) values ("test3", "aAあ亜", "$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2", "aAあ亜0123456789プロファイル");

INSERT INTO `grumbles` (`pk`, `content`, `user_id`, `created_at`) VALUES
('01FWMM3VMDVEG35G3ER9SG5G94', 'a', 'test1', '2022-02-24 10:04:21'),
('01FX1T8DPVG182CA9XJH392M47', 'l', 'test1', '2022-03-01 13:01:50'),
('01FX4VZHEN47XV7E3J6PFPDEJM', 'ああああああああああああ\naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n\nafgre\nrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr', 'test1', '2022-03-02 17:29:39'),
('01FX5AEM3549N0937KGBV34D9J', 'aaaaaaaaaaaaaaaaaaaaa\nrrrrrrrrrrrrrrrrr\ngbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb', 'test1', '2022-03-02 21:42:33'),
('01FX73CT6H9P77JZ9RWDZ9X2SC', 'gea a', 'test1', '2022-03-03 14:17:43'),
