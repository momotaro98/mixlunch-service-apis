
BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS `userblocklists` (
`blocker` CHAR (50) CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci` NOT NULL,
`blockee` CHAR (50) CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci` NOT NULL,
`createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updatedAt` DATETIME ON UPDATE CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`blocker`, `blockee`),
INDEX `userblocklists_ibfk_1` (`blocker`),
CONSTRAINT `userblocklists_ibfk_1` FOREIGN KEY (`blocker`) REFERENCES `users` (`userId`) ON DELETE CASCADE,
INDEX `userblocklists_ibfk_2` (`blockee`),
CONSTRAINT `userblocklists_ibfk_2` FOREIGN KEY (`blockee`) REFERENCES `users` (`userId`) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS `partymemberreviews` (
`partyId` INT (11) NOT NULL,
`reviewer` CHAR (50) CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci` NOT NULL,
`reviewee` CHAR (50) CHARACTER SET `utf8mb4` COLLATE `utf8mb4_general_ci` NOT NULL,
`score` DOUBLE NOT NULL,
`comments` TEXT CHARACTER SET `utf8mb4`,
`createdAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
`updatedAt` DATETIME ON UPDATE CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`partyId`, `reviewer`, `reviewee`),
INDEX `partymemberreviews_ibfk_1` (`partyId`),
CONSTRAINT `partymemberreviews_ibfk_1` FOREIGN KEY (`partyId`) REFERENCES `parties` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
INDEX `partymemberreviews_ibfk_2` (`reviewer`),
CONSTRAINT `partymemberreviews_ibfk_2` FOREIGN KEY (`reviewer`) REFERENCES `users` (`userId`) ON DELETE CASCADE,
INDEX `partymemberreviews_ibfk_3` (`reviewee`),
CONSTRAINT `partymemberreviews_ibfk_3` FOREIGN KEY (`reviewee`) REFERENCES `users` (`userId`) ON DELETE CASCADE
);

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;