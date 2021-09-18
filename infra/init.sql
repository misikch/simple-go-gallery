CREATE DATABASE IF NOT EXISTS `gallery`;
USE `gallery`;

CREATE TABLE IF NOT EXISTS `photos` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_id` int(11) NOT NULL,
    `path` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;