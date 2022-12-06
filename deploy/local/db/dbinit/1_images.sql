CREATE TABLE `photoshare`.`images` (
    id int(11) not null primary key auto_increment,
    image_name varchar(255) not null,
    tag varchar(50),
    title varchar(50),
    `description` varchar(255),
    date_added datetime not null
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;