package userservice

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/momotaro98/stew"
)

type UserFullQueryDto struct {
	userId             string
	name               string
	email              string
	nickName           sql.NullString
	sex                string
	birthday           time.Time
	photoUrl           sql.NullString
	latitude           float64
	longitude          float64
	positionName       sql.NullString
	academicBackground sql.NullString
	company            sql.NullString
	selfIntroduction   sql.NullString
	userlangs          []string
	useroccupations    []uint8
	usertags           []uint16
	blockingUsers      []string
}

type IUserQueryRepository interface {
	QueryUserFullByUsingUserId(userId string) (*UserFullQueryDto, error)
	QueryUserBlockWhereBlocker(blocker string) ([]*UserBlockQueryDto, error)
}

var _ IUserQueryRepository = (*realUserQueryRepository)(nil)

type realUserQueryRepository struct {
	db SqlDb
}

func ProvideUserQueryRepository(db SqlDb) IUserQueryRepository {
	return &realUserQueryRepository{
		db: db,
	}
}

func (r *realUserQueryRepository) QueryUserFullByUsingUserId(userId string) (*UserFullQueryDto, error) {
	queryAUser := `SELECT u.userId
			,u.name
			,u.email
			,u.nickName
			,u.sex
			,u.birthday
			,u.photoUrl
			,p.name AS positionName
			,u.academicBackground
			,u.company
			,u.selfIntroduction
		FROM users AS u
		LEFT JOIN positions AS p ON u.positionId=p.positionId
		WHERE u.userId = ?`
	var u UserFullQueryDto
	if err := r.db.QueryRow(queryAUser, userId).Scan(
		&u.userId, &u.name, &u.email, &u.nickName, &u.sex,
		&u.birthday, &u.photoUrl, &u.positionName,
		&u.academicBackground, &u.company, &u.selfIntroduction); err != nil {
		return nil, err
	}

	var err error

	// userlocations
	u.latitude, u.longitude, err = r.queryUserLocation(userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	// userlangs
	u.userlangs, err = r.queryUserLangs(userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	// useroccupations
	u.useroccupations, err = r.queryUserOccupations(userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	// usertags
	u.usertags, err = r.queryUserTags(userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	// blocking users
	userBlockedQDto, err := r.QueryUserBlockWhereBlocker(userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	for _, ub := range userBlockedQDto {
		u.blockingUsers = append(u.blockingUsers, ub.blockee)
	}

	return &u, nil
}

func (r *realUserQueryRepository) queryUserLocation(userId string) (latitude, longitude float64, err error) {
	var loc = struct {
		lat float64
		lng float64
	}{
		lat: latitude,
		lng: longitude,
	}
	if err := r.db.QueryRow(`
		SELECT latitude, longitude
		FROM userlocations
		WHERE userId = ?`, userId).Scan(
		&loc.lat, &loc.lng,
	); err != nil {
		return 0.0, 0.0, stew.Wrap(err)
	}
	return loc.lat, loc.lng, nil
}

func (r *realUserQueryRepository) queryUserLangs(userId string) (langs []string, err error) {
	rows, err := r.db.Query(`SELECT lang FROM userlangs
		WHERE userId = ?`, userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	for rows.Next() {
		var lang string
		if err := rows.Scan(&lang); err != nil {
			return nil, stew.Wrap(err)
		}
		langs = append(langs, lang)
	}
	return langs, nil
}

func (r *realUserQueryRepository) queryUserOccupations(userId string) (occupations []uint8, err error) {
	// [Note] For now, occupations table master is not needed for API since
	//        Front side has the occupation master instead.
	//rows, err := r.db.Query(`SELECT o.name AS occupationName
	//	FROM useroccupations AS uo
	//	INNER JOIN occupations AS o ON uo.occupationId=o.occupationId
	//	WHERE uo.userId = ?`, userId)
	rows, err := r.db.Query(`SELECT uo.occupationId AS occupationName
		FROM useroccupations AS uo
		WHERE uo.userId = ?`, userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	for rows.Next() {
		//var occupationName string
		var occupationID uint8
		if err := rows.Scan(&occupationID); err != nil {
			return nil, stew.Wrap(err)
		}
		occupations = append(occupations, occupationID)
	}
	return occupations, nil
}

func (r *realUserQueryRepository) queryUserTags(userId string) (tagIds []uint16, err error) {
	rows, err := r.db.Query(`SELECT tagId FROM usertags
		WHERE userId = ?`, userId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	for rows.Next() {
		var tagId uint16
		if err := rows.Scan(&tagId); err != nil {
			return nil, stew.Wrap(err)
		}
		tagIds = append(tagIds, tagId)
	}
	return tagIds, nil
}

type UserBlockQueryDto struct {
	blocker   string
	blockee   string
	createdAt time.Time
}

func (r *realUserQueryRepository) QueryUserBlockWhereBlocker(blocker string) ([]*UserBlockQueryDto, error) {
	rows, err := r.db.Query(`
		SELECT blocker, blockee, createdAt FROM userblocklists
		WHERE blocker = ?`, blocker)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	var ret []*UserBlockQueryDto
	for rows.Next() {
		var qDto UserBlockQueryDto
		if err := rows.Scan(&qDto.blocker, &qDto.blockee, &qDto.createdAt); err != nil {
			return nil, stew.Wrap(err)
		}
		ret = append(ret, &qDto)
	}
	return ret, nil
}

type UserCommandDto struct {
	userId             string
	name               string
	email              string
	nickName           sql.NullString
	sex                string
	birthday           time.Time
	photoUrl           sql.NullString
	latitude           float64
	longitude          float64
	positionId         sql.NullInt32
	academicBackground sql.NullString
	company            sql.NullString
	selfIntroduction   sql.NullString
	userlangs          []string
	occupationIDs      []uint8
	usertags           []uint16
}

type IUserCommandRepository interface {
	InsertUserInfo(user *UserCommandDto) error
	InsertUserBlock(userBlock *UserBlockCommandDto) error
}

var _ IUserCommandRepository = (*realUserCommandRepository)(nil)

type realUserCommandRepository struct {
	db SqlDb
}

func ProvideUserCommandRepository(db SqlDb) IUserCommandRepository {
	return &realUserCommandRepository{
		db: db,
	}
}

func (r *realUserCommandRepository) InsertUserInfo(u *UserCommandDto) error {
	tx, err := r.db.Begin()
	if err != nil {
		return stew.Wrap(err)
	}

	// users table
	_, err = tx.Exec(`
		INSERT INTO users (userId, name, email, nickName, sex, birthday, photoUrl, positionId, academicBackground, company, selfIntroduction)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, u.userId, u.name, u.email, u.nickName, u.sex, u.birthday, u.photoUrl, u.positionId, u.academicBackground, u.company, u.selfIntroduction)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		var e *mysql.MySQLError
		if errors.As(err, &e) {
			switch e.Number {
			case RepoErrCodeMapToRDBMS[DuplicatePrimaryKeyErrorCode]:
				return NewDuplicatePrimaryKeyError(err)
			}
		}
		return stew.Wrap(err)
	}

	// userlocations table
	_, err = tx.Exec(`
		INSERT INTO userlocations (userId, latitude, longitude)
		VALUES (?, ?, ?)
		`, u.userId, u.latitude, u.longitude)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		return stew.Wrap(err)
	}

	// userlangs table
	for _, l := range u.userlangs {
		_, err := tx.Exec(`
			INSERT INTO userlangs (userId, lang)
			VALUES (?, ?)
		`, u.userId, l)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				panic(err)
			}
			return err
		}
	}

	// useroccupations table
	for _, oID := range u.occupationIDs {
		_, err := tx.Exec(`
			INSERT INTO useroccupations (userId, occupationId)
			VALUES (?, ?)
        `, u.userId, oID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				panic(err)
			}
			return stew.Wrap(err)
		}
	}

	// usertags table
	for _, tagId := range u.usertags {
		_, err := tx.Exec(`
			INSERT INTO usertags (userId, tagId)
			VALUES (?, ?)
		`, u.userId, tagId)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				panic(err)
			}
			return stew.Wrap(err)
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		return stew.Wrap(err)
	}
	return nil
}

type UserBlockCommandDto struct {
	blocker string
	blockee string
}

func (r *realUserCommandRepository) InsertUserBlock(ub *UserBlockCommandDto) error {
	tx, err := r.db.Begin()
	if err != nil {
		return stew.Wrap(err)
	}

	_, err = tx.Exec(`
		INSERT INTO userblocklists (blocker, blockee)
		VALUES (?, ?)
		`, ub.blocker, ub.blockee)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return stew.Wrap(err)
		}
		var e *mysql.MySQLError
		if errors.As(err, &e) {
			switch e.Number {
			case RepoErrCodeMapToRDBMS[DuplicatePrimaryKeyErrorCode]:
				return NewDuplicatePrimaryKeyError(err)
			case RepoErrCodeMapToRDBMS[NoReferenceRowErrorCode]:
				return NewNoReferenceRowError(err)
			}
		}
		return stew.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
		return stew.Wrap(err)
	}

	return nil
}
