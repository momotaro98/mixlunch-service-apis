
BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE `locations`;

ALTER TABLE `users` DROP FOREIGN KEY `users_ibfk_1`;
ALTER TABLE `users` DROP INDEX `users_ibfk_1`;
ALTER TABLE `users` DROP COLUMN `locationId`;

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;
