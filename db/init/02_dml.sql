SET CHARACTER_SET_CLIENT = utf8mb4;
SET CHARACTER_SET_CONNECTION = utf8mb4;
-- category master
-- categories table
INSERT INTO categories (categoryId, name) VALUES (0, 'root');
INSERT INTO categories (categoryId, name) VALUES (1, 'IT・情報');
INSERT INTO categories (categoryId, name) VALUES (2, 'スポーツ');
INSERT INTO categories (categoryId, name) VALUES (3, 'グルメ');
INSERT INTO categories (categoryId, name) VALUES (4, 'プログラミング');
INSERT INTO categories (categoryId, name) VALUES (5, 'AI');
-- categoriestree table
INSERT INTO categoriestree (parentId, childId) VALUES (0, 1);
INSERT INTO categoriestree (parentId, childId) VALUES (0, 2);
INSERT INTO categoriestree (parentId, childId) VALUES (0, 3);
INSERT INTO categoriestree (parentId, childId) VALUES (1, 4);
INSERT INTO categoriestree (parentId, childId) VALUES (1, 5);
-- tag master
-- tagtypes table
INSERT INTO tagtypes (tagTypeId, name) VALUES (1, 'Interest');
INSERT INTO tagtypes (tagTypeId, name) VALUES (2, 'Skill');
-- tags table
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('Python', 2, 4);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('Vue.js', 2, 4);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('Ruby on Rails', 2, 4);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('TensorFlow', 2, 5);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('ディープラーニング', 2, 5);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('野球', 1, 2);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('サッカー', 1, 2);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('焼き肉', 1, 3);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('寿司', 1, 3);
INSERT INTO tags (name, tagTypeId, categoryId) VALUES ('ハンバーガー', 1, 3);

-- positions
INSERT INTO positions (name) VALUES ('Employee');
SET @firstPositionId=LAST_INSERT_ID();
INSERT INTO positions (name) VALUES ('Manager');
SET @secondPositionId=LAST_INSERT_ID();

-- occupations
INSERT INTO occupations (name) VALUES ('Application Engineer');
SET @firstOccupationId=LAST_INSERT_ID();
INSERT INTO occupations (name) VALUES ('Data Scientist');
SET @secondOccupationId=LAST_INSERT_ID();
INSERT INTO occupations (name) VALUES ('Producer');
SET @thirdOccupationId=LAST_INSERT_ID();
INSERT INTO occupations (name) VALUES ('Dummy'); -- 4
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 10
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 20
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 30
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 40
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 50
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 60
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 70
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 80
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 90
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy');
INSERT INTO occupations (name) VALUES ('Dummy'); -- 100

-- locationtypes
INSERT INTO locationtypes (id, name) VALUES (0, 'Geographic');
INSERT INTO locationtypes (id, name) VALUES (1, 'Online');

-- users
SET @userId1:='xDlXdTXw5eV7jC7ETxX59gUk71J2';
SET @userId2:='plpC2nFJ5ifFN5WaFsjhIce1KFy2';
SET @userId3:='agKZdfhGOQMPQcNuCrR3xyfrrku1';
SET @userId4:='tVI1i49UBFPvnse9NuRR8m4X6342';
INSERT INTO users (userId, name, email, nickName, sex, birthday, positionId, academicBackground, company, selfIntroduction)
VALUES (@userId1, 'Shintaro Ikeda', 'ikenshirogivenup98@gmail.com', 'Shin', '1', '1991-08-04', @firstPositionId, '早稲田大学', 'リクルート', 'こんにちは。よろしくです。');
INSERT INTO users (userId, name, email, nickName, sex, birthday, positionId)
VALUES (@userId2, 'Yota Ishikawa', 'yota.ishikawa@rakuten.com', 'Yota', '1', '1993-03-04', @firstPositionId);
INSERT INTO users (userId, name, email, nickName, sex, birthday, positionId)
VALUES (@userId3, 'John X', 'dev.momotaro98@gmail.com', 'John', '1', '1995-03-01', NULL);
INSERT INTO users (userId, name, email, nickName, sex, birthday, positionId)
VALUES (@userId4, 'Shinjiro Yamada', 'ikenshirodiscovery@yahoo.co.jp', 'Shin', '1', '1997-02-28', @secondPositionId);

-- userlocations
-- userId1
INSERT INTO userlocations (userId, latitude, longitude) VALUES (@userId1, 35.681236, 139.767125); -- 東京駅
-- userId2
INSERT INTO userlocations (userId, latitude, longitude) VALUES (@userId2, 35.658034, 139.701636); -- 渋谷駅
-- userId3
INSERT INTO userlocations (userId, latitude, longitude) VALUES (@userId3, 35.689607, 139.700571); -- 新宿駅
-- userId4
INSERT INTO userlocations (userId, latitude, longitude) VALUES (@userId4, 34.702485, 135.495951); -- 大阪駅

-- userlangs
-- Based on ISO 639-1
-- userId1
INSERT INTO userlangs (userId, lang) VALUES (@userId1, 'ja');
INSERT INTO userlangs (userId, lang) VALUES (@userId1, 'en');
-- userId2
INSERT INTO userlangs (userId, lang) VALUES (@userId2, 'ja');
INSERT INTO userlangs (userId, lang) VALUES (@userId2, 'en');
-- userId3
INSERT INTO userlangs (userId, lang) VALUES (@userId3, 'ja');
INSERT INTO userlangs (userId, lang) VALUES (@userId3, 'en');
INSERT INTO userlangs (userId, lang) VALUES (@userId3, 'fr');
INSERT INTO userlangs (userId, lang) VALUES (@userId3, 'zh');
-- userId4
INSERT INTO userlangs (userId, lang) VALUES (@userId4, 'ja');

-- useroccupations
-- userId1
INSERT INTO useroccupations (userId, occupationId) VALUES (@userId1, @firstOccupationId);
INSERT INTO useroccupations (userId, occupationId) VALUES (@userId1, @secondOccupationId);
-- userId2
INSERT INTO useroccupations (userId, occupationId) VALUES (@userId2, @firstOccupationId);
INSERT INTO useroccupations (userId, occupationId) VALUES (@userId2, @secondOccupationId);
-- userId3
INSERT INTO useroccupations (userId, occupationId) VALUES (@userId3, @firstOccupationId);
INSERT INTO useroccupations (userId, occupationId) VALUES (@userId3, @secondOccupationId);
INSERT INTO useroccupations (userId, occupationId) VALUES (@userId3, @thirdOccupationId);
-- userId4
INSERT INTO useroccupations (userId, occupationId) VALUES (@userId4, @thirdOccupationId);

-- usertags
-- userId1
INSERT INTO usertags (userId, tagId) VALUES (@userId1, 1);
INSERT INTO usertags (userId, tagId) VALUES (@userId1, 2);
INSERT INTO usertags (userId, tagId) VALUES (@userId1, 3);
INSERT INTO usertags (userId, tagId) VALUES (@userId1, 6);
INSERT INTO usertags (userId, tagId) VALUES (@userId1, 10);
-- userId2
INSERT INTO usertags (userId, tagId) VALUES (@userId2, 1);
INSERT INTO usertags (userId, tagId) VALUES (@userId2, 2);
INSERT INTO usertags (userId, tagId) VALUES (@userId2, 4);
INSERT INTO usertags (userId, tagId) VALUES (@userId2, 7);
INSERT INTO usertags (userId, tagId) VALUES (@userId2, 8);
-- userId3
INSERT INTO usertags (userId, tagId) VALUES (@userId3, 4);
INSERT INTO usertags (userId, tagId) VALUES (@userId3, 5);
-- userId4
INSERT INTO usertags (userId, tagId) VALUES (@userId4, 9);
INSERT INTO usertags (userId, tagId) VALUES (@userId4, 10);

-- userschedules
-- userId1
INSERT INTO userschedules (userId, fromDateTime, toDateTime, locationTypeId) VALUES (@userId1, '2018-11-01 3:00:00', '2018-11-01 5:00:00', 0);
SET @usID1=LAST_INSERT_ID();
-- userId2
INSERT INTO userschedules (userId, fromDateTime, toDateTime, locationTypeId) VALUES (@userId2, '2018-11-01 3:00:00', '2018-11-01 5:00:00', 0);
SET @usID2=LAST_INSERT_ID();
-- userId3
INSERT INTO userschedules (userId, fromDateTime, toDateTime, locationTypeId) VALUES (@userId3, '2018-11-01 3:00:00', '2018-11-01 5:00:00', 0);
SET @usID3=LAST_INSERT_ID();
-- userId4
INSERT INTO userschedules (userId, fromDateTime, toDateTime, locationTypeId) VALUES (@userId4, '2018-11-01 3:00:00', '2018-11-01 5:00:00', 0);
SET @usID4=LAST_INSERT_ID();
-- userId5 (Online)
INSERT INTO userschedules (userId, fromDateTime, toDateTime, locationTypeId) VALUES (@userId4, '2018-11-01 3:00:00', '2018-11-01 5:00:00', 1);
SET @usID5=LAST_INSERT_ID();

-- userscheduletags
-- userId1
INSERT INTO userscheduletags (userScheduleId, tagId) VALUES (@usID1, 1);
INSERT INTO userscheduletags (userScheduleId, tagId) VALUES (@usID1, 2);
-- userId2
INSERT INTO userscheduletags (userScheduleId, tagId) VALUES (@usID2, 1);
INSERT INTO userscheduletags (userScheduleId, tagId) VALUES (@usID2, 3);
-- userId3
-- userId4
INSERT INTO userscheduletags (userScheduleId, tagId) VALUES (@usID4, 1);
INSERT INTO userscheduletags (userScheduleId, tagId) VALUES (@usID4, 3);
INSERT INTO userscheduletags (userScheduleId, tagId) VALUES (@usID4, 4);
-- userId5

-- userschedulelocations
-- userId1
INSERT INTO userschedulelocations (userScheduleId, latitude, longitude) VALUES (@usID1, 35.681236, 139.767125); -- 東京駅
-- userId2
INSERT INTO userschedulelocations (userScheduleId, latitude, longitude) VALUES (@usID2, 35.681236, 139.767125); -- 東京駅
-- userId3
INSERT INTO userschedulelocations (userScheduleId, latitude, longitude) VALUES (@usID3, 35.681236, 139.767125); -- 東京駅
-- userId4
INSERT INTO userschedulelocations (userScheduleId, latitude, longitude) VALUES (@usID4, 35.681236, 139.767125); -- 東京駅
-- userId5 doesn't have userschedulelocations since it's online type

-- parties
INSERT INTO parties (startFrom, endTo)
VALUES ('2018-11-01 3:00:00', '2018-11-01 4:00:00');
SET @lastPartyID=LAST_INSERT_ID();

-- partymembers
INSERT INTO partymembers (partyId, userId)
VALUES (@lastPartyID, @userId1);
INSERT INTO partymembers (partyId, userId)
VALUES (@lastPartyID, @userId3);

-- partytags
INSERT INTO partytags (partyId, tagId)
VALUES (@lastPartyID, 2);
INSERT INTO partytags (partyId, tagId)
VALUES (@lastPartyID, 3);

-- partymemberreviews
INSERT INTO partymemberreviews (partyId, reviewer, reviewee, score, comments)
VALUES (@lastPartyID, @userId1, @userId3, 0.8, 'He was good.');

-- userblocklists
INSERT INTO userblocklists (blocker, blockee)
VALUES (@userId1, @userId4);
INSERT INTO userblocklists (blocker, blockee)
VALUES (@userId2, @userId4);
INSERT INTO userblocklists (blocker, blockee)
VALUES (@userId4, @userId1);
