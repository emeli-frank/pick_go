CREATE DATABASE pick CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'pick'@'localhost' IDENTIFIED BY 'pick';
GRANT ALL ON pick.* TO 'pick'@'localhost';
-- USE pick
