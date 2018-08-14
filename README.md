# TS_users
### Tools: GoLang, MySQL 
APIs to:
  1. Create users
  2. Declare a random person (of all the users created in the last 30 mins) as a winner

# Database Table
CREATE TABLE `TS_users` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Primary key users table',
  `first_name` varbinary(2550) DEFAULT NULL,
  `last_name` varbinary(2550) DEFAULT NULL,
  `email` varbinary(2550) DEFAULT NULL,
  `phone` varchar(12) DEFAULT NULL,
  `contact_me` tinyint(1) DEFAULT NULL,
  `how_long` varchar(12) DEFAULT NULL,
  `created_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record created time',
  `modified_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record updated time',
  `is_winner` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Postman API:
https://www.getpostman.com/collections/f00393b3e92571009b8b
