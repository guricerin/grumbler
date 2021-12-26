use `grumbler_db`;

-- pass: aaaaaaaa
insert into `users` (`id`, `name`, `password`, `profile`) values ("test1", "tEsT", "$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2", "test profile");
insert into `users` (`id`, `name`, `password`, `profile`) values ("test2", "尾骶骨テスト", "$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2", " つちよし𠮷test ");
insert into `users` (`id`, `name`, `password`, `profile`) values ("test3", "aAあ亜", "$2a$10$I90U5V.7xG34M6Q/NYkt1OWZ0oz19BSLNAZcVrZBbDcY0coXXSLH2", "aAあ亜0123456789プロファイル");
