CREATE TABLE `photoshare`.`images` (
    id int(11) not null primary key auto_increment,
    uri varchar(255) not null,
    title varchar(100) not null,
    tag varchar(100),
    date_added datetime not null
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;