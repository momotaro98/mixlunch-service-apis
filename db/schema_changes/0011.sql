
BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `userschedulelocations` (
`userScheduleId` INT (11) NOT NULL,
`latitude` DOUBLE NOT NULL,
`longitude` DOUBLE NOT NULL,
`createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updatedAt` DATETIME ON UPDATE CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`userScheduleId`),
INDEX `userschedulelocations_ibfk_1` (`userScheduleId`),
CONSTRAINT `userschedulelocations_ibfk_1` FOREIGN KEY (`userScheduleId`) REFERENCES `userschedules` (`userScheduleId`) ON DELETE CASCADE ON UPDATE CASCADE
);

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;