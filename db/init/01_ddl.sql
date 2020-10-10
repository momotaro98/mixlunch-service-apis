CREATE TABLE IF NOT EXISTS positions (
    positionId SMALLINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50),
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (positionId)
);

CREATE TABLE IF NOT EXISTS occupations (
    occupationId SMALLINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50),
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (occupationId)
);

CREATE TABLE IF NOT EXISTS locationtypes (
    id TINYINT NOT NULL,
    name VARCHAR(50) CHARACTER SET utf8 NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users (
    userId CHAR(50) NOT NULL,
    name VARCHAR(200) CHARACTER SET utf8 NOT NULL,
    email VARCHAR(256) NOT NULL,
    nickName VARCHAR(50) CHARACTER SET utf8,
    sex CHAR(1) NOT NULL DEFAULT '0',
    birthday DATE NOT NULL,
    photoUrl TEXT CHARACTER SET utf8,
    positionId SMALLINT,
    academicBackground TEXT CHARACTER SET utf8mb4,
    company TEXT CHARACTER SET utf8mb4,
    selfIntroduction TEXT CHARACTER SET utf8mb4,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userId),
    CONSTRAINT users_ibfk_3 FOREIGN KEY(positionId) REFERENCES positions(positionId) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS userlocations (
    userId CHAR(50) NOT NULL,
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userId),
    CONSTRAINT userlocations_ibfk_1 FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS userlangs (
    userId CHAR(50) NOT NULL,
    lang CHAR(2) NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userId, lang),
    CONSTRAINT userlangs_ibfk_1 FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS useroccupations (
    userId CHAR(50) NOT NULL,
    occupationId SMALLINT NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userId, occupationId),
    CONSTRAINT useroccupations_ibfk_1 FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE,
    CONSTRAINT useroccupations_ibfk_2 FOREIGN KEY(occupationId) REFERENCES occupations(occupationId) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
    categoryId SMALLINT NOT NULL,
    name VARCHAR(50) CHARACTER SET utf8 NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (categoryId)
);

CREATE TABLE IF NOT EXISTS categoriestree (
    parentId SMALLINT NOT NULL,
    childId SMALLINT NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (parentId, childId),
    CONSTRAINT categoriestree_ibfk_1 FOREIGN KEY(parentId) REFERENCES categories(categoryId) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT categoriestree_ibfk_2 FOREIGN KEY(childId) REFERENCES categories(categoryId) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS tagtypes (
    tagTypeId TINYINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (tagTypeId)
);

CREATE TABLE IF NOT EXISTS tags (
    tagId MEDIUMINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) CHARACTER SET utf8 NOT NULL,
    tagTypeId TINYINT NOT NULL,
    categoryId SMALLINT NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (tagId),
    CONSTRAINT tags_ibfk_1 FOREIGN KEY(tagTypeId) REFERENCES tagtypes(tagTypeId) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT tags_ibfk_2 FOREIGN KEY(categoryId) REFERENCES categories(categoryId) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS usertags (
    userId CHAR(50) NOT NULL,
    tagId MEDIUMINT NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userId, tagId),
    CONSTRAINT usertags_ibfk_1 FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE,
    CONSTRAINT usertags_ibfk_2 FOREIGN KEY(tagId) REFERENCES tags(tagId) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS userschedules (
    userScheduleId INT NOT NULL AUTO_INCREMENT,
    userId CHAR(50) NOT NULL,
    fromDateTime DATETIME NOT NULL,
    toDateTime DATETIME NOT NULL,
    locationTypeId TINYINT NOT NULL DEFAULT 0,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userScheduleId),
    CONSTRAINT userschedules_ibfk_1 FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE,
    CONSTRAINT userschedules_ibfk_2 FOREIGN KEY(locationTypeId) REFERENCES locationtypes(id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS userscheduletags (
    userScheduleId INT NOT NULL,
    tagId MEDIUMINT NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userScheduleId, tagId),
    CONSTRAINT userscheduletags_ibfk_1 FOREIGN KEY(userScheduleId) REFERENCES userschedules(userScheduleId) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT userscheduletags_ibfk_2 FOREIGN KEY(tagId) REFERENCES tags(tagId) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS userschedulelocations (
    userScheduleId INT NOT NULL,
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (userScheduleId),
    CONSTRAINT userschedulelocations_ibfk_1 FOREIGN KEY(userScheduleId) REFERENCES userschedules(userScheduleId) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS parties (
    id INT NOT NULL AUTO_INCREMENT,
    startFrom DATETIME NOT NULL,
    endTo DATETIME NOT NULL,
    chatRoomId CHAR(50),
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS partymembers (
    partyMemberId INT NOT NULL AUTO_INCREMENT,
    partyId INT NOT NULL,
    userId CHAR(50) NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT party_user UNIQUE (partyId, userId),
    PRIMARY KEY (partyMemberId),
    CONSTRAINT partymembers_ibfk_1 FOREIGN KEY(partyId) REFERENCES parties(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT partymembers_ibfk_2 FOREIGN KEY(userId) REFERENCES users(userId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS partytags (
    partyId INT NOT NULL,
    tagId MEDIUMINT NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (partyId, tagId),
    CONSTRAINT partytags_ibfk_1 FOREIGN KEY(partyId) REFERENCES parties(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT partytags_ibfk_2 FOREIGN KEY(tagId) REFERENCES tags(tagId) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS userblocklists (
    blocker CHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    blockee CHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (blocker, blockee),
    CONSTRAINT userblocklists_ibfk_1 FOREIGN KEY(blocker) REFERENCES users(userId) ON DELETE CASCADE,
    CONSTRAINT userblocklists_ibfk_2 FOREIGN KEY(blockee) REFERENCES users(userId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS partymemberreviews (
    partyId INT NOT NULL,
    reviewer CHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    reviewee CHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    score DOUBLE NOT NULL,
    comments TEXT CHARACTER SET utf8mb4,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (partyId, reviewer, reviewee),
    CONSTRAINT partymemberreviews_ibfk_1 FOREIGN KEY(partyId) REFERENCES parties(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT partymemberreviews_ibfk_2 FOREIGN KEY(reviewer) REFERENCES users(userId) ON DELETE CASCADE,
    CONSTRAINT partymemberreviews_ibfk_3 FOREIGN KEY(reviewee) REFERENCES users(userId) ON DELETE CASCADE
);