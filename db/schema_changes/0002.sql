BEGIN;

SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE `userschedules` CHANGE COLUMN `id` `id` int(10) unsigned NOT NULL; -- AUTO_INCREMENTなカラムはプライマリーキーMUSTなのでAUTO_INCREMENTじゃなくさせる
ALTER TABLE `userschedules` DROP INDEX `PRIMARY`;
ALTER TABLE `userschedules` DROP COLUMN `id`;
ALTER TABLE `userschedules` ADD COLUMN `userScheduleId` INT (11) NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST;

ALTER TABLE `partymembers` CHANGE COLUMN `id` `id` int(10) unsigned NOT NULL; -- AUTO_INCREMENTなカラムはプライマリーキーMUSTなのでAUTO_INCREMENTじゃなくさせる
ALTER TABLE `partymembers` DROP INDEX `PRIMARY`;
ALTER TABLE `partymembers` DROP COLUMN `id`;
ALTER TABLE `partymembers` ADD COLUMN `partyMemberId` INT (11) NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST;

ALTER TABLE `partymembers` DROP FOREIGN KEY `partymembers_ibfk_1`;

ALTER TABLE `parties` CHANGE COLUMN `id` `id` INT (11) NOT NULL AUTO_INCREMENT;
ALTER TABLE `partymembers` CHANGE COLUMN `partyId` `partyId` INT (11) NOT NULL;
ALTER TABLE `partymembers` ADD CONSTRAINT `party_user` UNIQUE INDEX (`partyId`, `userId`);

ALTER TABLE `partymembers` ADD CONSTRAINT `partymembers_ibfk_1` FOREIGN KEY (partyId) REFERENCES parties(id);

SET FOREIGN_KEY_CHECKS = 1;

COMMIT;
