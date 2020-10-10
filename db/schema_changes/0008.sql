
BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `userlangs` DROP FOREIGN KEY `userlangs_ibfk_1`;
ALTER TABLE `userlangs` ADD CONSTRAINT `userlangs_ibfk_1` FOREIGN KEY (`userId`) REFERENCES `users` (`userId`) ON DELETE CASCADE;
ALTER TABLE `useroccupations` DROP FOREIGN KEY `useroccupations_ibfk_1`;
ALTER TABLE `useroccupations` DROP FOREIGN KEY `useroccupations_ibfk_2`;
ALTER TABLE `useroccupations` ADD CONSTRAINT `useroccupations_ibfk_1` FOREIGN KEY (`userId`) REFERENCES `users` (`userId`) ON DELETE CASCADE;
ALTER TABLE `useroccupations` ADD CONSTRAINT `useroccupations_ibfk_2` FOREIGN KEY (`occupationId`) REFERENCES `occupations` (`occupationId`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `tags` DROP FOREIGN KEY `tags_ibfk_1`;
ALTER TABLE `tags` DROP FOREIGN KEY `tags_ibfk_2`;
ALTER TABLE `tags` ADD CONSTRAINT `tags_ibfk_1` FOREIGN KEY (`tagTypeId`) REFERENCES `tagtypes` (`tagTypeId`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `tags` ADD CONSTRAINT `tags_ibfk_2` FOREIGN KEY (`categoryId`) REFERENCES `categories` (`categoryId`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `users` DROP FOREIGN KEY `users_ibfk_3`;
ALTER TABLE `users` DROP FOREIGN KEY `users_ibfk_1`;
ALTER TABLE `users` ADD CONSTRAINT `users_ibfk_1` FOREIGN KEY (`locationId`) REFERENCES `locations` (`locationId`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `users` ADD CONSTRAINT `users_ibfk_3` FOREIGN KEY (`positionId`) REFERENCES `positions` (`positionId`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `categoriestree` DROP FOREIGN KEY `categoriestree_ibfk_1`;
ALTER TABLE `categoriestree` DROP FOREIGN KEY `categoriestree_ibfk_2`;
ALTER TABLE `categoriestree` ADD CONSTRAINT `categoriestree_ibfk_1` FOREIGN KEY (`parentId`) REFERENCES `categories` (`categoryId`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `categoriestree` ADD CONSTRAINT `categoriestree_ibfk_2` FOREIGN KEY (`childId`) REFERENCES `categories` (`categoryId`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `partytags` DROP FOREIGN KEY `partytags_ibfk_1`;
ALTER TABLE `partytags` DROP FOREIGN KEY `partytags_ibfk_2`;
ALTER TABLE `partytags` ADD CONSTRAINT `partytags_ibfk_2` FOREIGN KEY (`tagId`) REFERENCES `tags` (`tagId`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `partytags` ADD CONSTRAINT `partytags_ibfk_1` FOREIGN KEY (`partyId`) REFERENCES `parties` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `userscheduletags` DROP FOREIGN KEY `userscheduletags_ibfk_2`;
ALTER TABLE `userscheduletags` DROP FOREIGN KEY `userscheduletags_ibfk_1`;
ALTER TABLE `userscheduletags` ADD CONSTRAINT `userscheduletags_ibfk_2` FOREIGN KEY (`tagId`) REFERENCES `tags` (`tagId`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `userscheduletags` ADD CONSTRAINT `userscheduletags_ibfk_1` FOREIGN KEY (`userScheduleId`) REFERENCES `userschedules` (`userScheduleId`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `usertags` DROP FOREIGN KEY `usertags_ibfk_1`;
ALTER TABLE `usertags` DROP FOREIGN KEY `usertags_ibfk_2`;
ALTER TABLE `usertags` ADD CONSTRAINT `usertags_ibfk_1` FOREIGN KEY (`userId`) REFERENCES `users` (`userId`) ON DELETE CASCADE;
ALTER TABLE `usertags` ADD CONSTRAINT `usertags_ibfk_2` FOREIGN KEY (`tagId`) REFERENCES `tags` (`tagId`) ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE `userschedules` DROP FOREIGN KEY `userschedules_ibfk_1`;
ALTER TABLE `userschedules` ADD CONSTRAINT `userschedules_ibfk_1` FOREIGN KEY (`userId`) REFERENCES `users` (`userId`) ON DELETE CASCADE;
ALTER TABLE `partymembers` DROP FOREIGN KEY `partymembers_ibfk_1`;
ALTER TABLE `partymembers` DROP FOREIGN KEY `partymembers_ibfk_2`;
ALTER TABLE `partymembers` ADD CONSTRAINT `partymembers_ibfk_1` FOREIGN KEY (`partyId`) REFERENCES `parties` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE `partymembers` ADD CONSTRAINT `partymembers_ibfk_2` FOREIGN KEY (`userId`) REFERENCES `users` (`userId`) ON DELETE CASCADE;

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;