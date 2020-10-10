
BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `userschedules` ADD FOREIGN KEY (`userId`) REFERENCES `users` (`userId`);
ALTER TABLE `partymembers` ADD FOREIGN KEY (`userId`) REFERENCES `users` (`userId`);

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;