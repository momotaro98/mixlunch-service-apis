BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `categoriestree` (
`parentId` SMALLINT (6) NOT NULL,
`childId` SMALLINT (6) NOT NULL,
PRIMARY KEY (`parentId`, `childId`),
FOREIGN KEY (`parentId`) REFERENCES `categories` (`categoryId`),
FOREIGN KEY (`childId`) REFERENCES `categories` (`categoryId`)
);
CREATE TABLE IF NOT EXISTS `departments` (
`departmentId` SMALLINT (6) NOT NULL AUTO_INCREMENT,
`name` VARCHAR (50) DEFAULT NULL,
PRIMARY KEY (`departmentId`)
);
CREATE TABLE IF NOT EXISTS `categories` (
`categoryId` SMALLINT (6) NOT NULL,
`name` VARCHAR (50) CHARACTER SET `utf8` NOT NULL,
PRIMARY KEY (`categoryId`)
);
CREATE TABLE IF NOT EXISTS `partytags` (
`partyId` INT (11) NOT NULL,
`tagId` MEDIUMINT (9) NOT NULL,
PRIMARY KEY (`partyId`, `tagId`),
FOREIGN KEY (`partyId`) REFERENCES `parties` (`id`),
FOREIGN KEY (`tagId`) REFERENCES `tags` (`tagId`)
);
CREATE TABLE IF NOT EXISTS `userlangs` (
`userId` CHAR (50) NOT NULL,
`lang` CHAR (2) NOT NULL,
PRIMARY KEY (`userId`, `lang`),
FOREIGN KEY (`userId`) REFERENCES `users` (`userId`)
);
CREATE TABLE IF NOT EXISTS `locations` (
`locationId` SMALLINT (6) NOT NULL AUTO_INCREMENT,
`name` VARCHAR (50) DEFAULT NULL,
PRIMARY KEY (`locationId`)
);
CREATE TABLE IF NOT EXISTS `tags` (
`tagId` MEDIUMINT (9) NOT NULL AUTO_INCREMENT,
`name` VARCHAR (50) CHARACTER SET `utf8` NOT NULL,
`tagTypeId` TINYINT (4) NOT NULL,
`categoryId` SMALLINT (6) NOT NULL,
PRIMARY KEY (`tagId`),
FOREIGN KEY (`tagTypeId`) REFERENCES `tagtypes` (`tagTypeId`),
FOREIGN KEY (`categoryId`) REFERENCES `categories` (`categoryId`)
);
CREATE TABLE IF NOT EXISTS `occupations` (
`occupationId` SMALLINT (6) NOT NULL AUTO_INCREMENT,
`name` VARCHAR (50) DEFAULT NULL,
PRIMARY KEY (`occupationId`)
);
CREATE TABLE IF NOT EXISTS `usertags` (
`userId` CHAR (50) NOT NULL,
`tagId` MEDIUMINT (9) NOT NULL,
PRIMARY KEY (`userId`, `tagId`),
FOREIGN KEY (`userId`) REFERENCES `users` (`userId`),
FOREIGN KEY (`tagId`) REFERENCES `tags` (`tagId`)
);
CREATE TABLE IF NOT EXISTS `positions` (
`positionId` SMALLINT (6) NOT NULL AUTO_INCREMENT,
`name` VARCHAR (50) DEFAULT NULL,
PRIMARY KEY (`positionId`)
);
CREATE TABLE IF NOT EXISTS `tagtypes` (
`tagTypeId` TINYINT (4) NOT NULL,
`name` VARCHAR (50) NOT NULL,
PRIMARY KEY (`tagTypeId`)
);
CREATE TABLE IF NOT EXISTS `users` (
`userId` CHAR (50) NOT NULL,
`name` VARCHAR (50) CHARACTER SET `utf8` NOT NULL,
`email` VARCHAR (256) DEFAULT NULL,
`nickName` VARCHAR (50) CHARACTER SET `utf8` DEFAULT NULL,
`sex` CHAR (1) NOT NULL DEFAULT '0',
`birthday` DATE DEFAULT NULL,
`locationId` SMALLINT (6) NOT NULL,
`departmentId` SMALLINT (6) DEFAULT NULL,
`positionId` SMALLINT (6) DEFAULT NULL,
`occupationId` SMALLINT (6) DEFAULT NULL,
PRIMARY KEY (`userId`),
FOREIGN KEY (`locationId`) REFERENCES `locations` (`locationId`),
FOREIGN KEY (`departmentId`) REFERENCES `departments` (`departmentId`),
FOREIGN KEY (`positionId`) REFERENCES `positions` (`positionId`),
FOREIGN KEY (`occupationId`) REFERENCES `occupations` (`occupationId`)
);
CREATE TABLE IF NOT EXISTS `userscheduletags` (
`userScheduleId` INT (11) NOT NULL,
`tagId` MEDIUMINT (9) NOT NULL,
PRIMARY KEY (`userScheduleId`, `tagId`),
FOREIGN KEY (`userScheduleId`) REFERENCES `userschedules` (`userScheduleId`),
FOREIGN KEY (`tagId`) REFERENCES `tags` (`tagId`)
);

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;
