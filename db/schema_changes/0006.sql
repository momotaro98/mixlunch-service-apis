BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `useroccupations` (
`userId` CHAR (50) NOT NULL,
`occupationId` SMALLINT (6) NOT NULL,
PRIMARY KEY (`userId`, `occupationId`),
FOREIGN KEY (`userId`) REFERENCES `users` (`userId`),
FOREIGN KEY (`occupationId`) REFERENCES `occupations` (`occupationId`)
);

ALTER TABLE `users` DROP FOREIGN KEY `users_ibfk_2`;
ALTER TABLE `users` DROP FOREIGN KEY `users_ibfk_4`;

ALTER TABLE `users` DROP COLUMN `occupationId`;
ALTER TABLE `users` DROP COLUMN `departmentId`;
ALTER TABLE `users` ADD COLUMN `academicBackground` TEXT CHARACTER SET `utf8mb4` AFTER `positionId`;
ALTER TABLE `users` ADD COLUMN `company` TEXT CHARACTER SET `utf8mb4` AFTER `academicBackground`;
ALTER TABLE `users` ADD COLUMN `selfIntroduction` TEXT CHARACTER SET `utf8mb4` AFTER `company`;
ALTER TABLE `users` CHANGE COLUMN `birthday` `birthday` DATE NOT NULL;
ALTER TABLE `users` CHANGE COLUMN `email` `email` VARCHAR (256) NOT NULL;

DROP TABLE `departments`;

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;
