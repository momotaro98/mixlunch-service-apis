
BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `locationtypes` (
`id` TINYINT (4) NOT NULL,
`name` VARCHAR (50) CHARACTER SET `utf8` NOT NULL,
`createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updatedAt` DATETIME ON UPDATE CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`id`)
);

INSERT INTO locationtypes (id, name) VALUES (0, 'Geographic');
INSERT INTO locationtypes (id, name) VALUES (1, 'Online');

ALTER TABLE `userschedules` ADD COLUMN `locationTypeId` TINYINT (4) NOT NULL DEFAULT 0 AFTER `toDateTime`;
ALTER TABLE `userschedules` ADD INDEX `userschedules_ibfk_2` (`locationTypeId`);
ALTER TABLE `userschedules` ADD CONSTRAINT `userschedules_ibfk_2` FOREIGN KEY (`locationTypeId`) REFERENCES `locationtypes` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;
