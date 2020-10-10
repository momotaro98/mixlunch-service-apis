package userscheduleservice

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/momotaro98/stew"
)

// UserScheduleDto is a data transfer object for userschedules table
type UserScheduleDto struct {
	userScheduleId int64
	userId         string
	fromDateTime   time.Time
	toDateTime     time.Time
	tagIds         []uint16
	latitude       float64
	longitude      float64
	locationTypeID int8
}

func newUserScheduleDtoForQuery(
	userScheduleId int64, userId string,
	fromDatetime, toDateTime time.Time,
	tagIds []uint16,
	latitude, longitude float64,
	locationTypeID int8,
) *UserScheduleDto {
	return &UserScheduleDto{
		userScheduleId: userScheduleId,
		userId:         userId,
		fromDateTime:   fromDatetime,
		toDateTime:     toDateTime,
		tagIds:         tagIds,
		latitude:       latitude,
		longitude:      longitude,
		locationTypeID: locationTypeID,
	}
}

// NewUserScheduleDtoForCommand generates UserScheduleDto instance
func NewUserScheduleDtoForCommand(
	userId string,
	fromDatetime, toDateTime time.Time,
	tagIds []uint16,
	latitude, longitude float64,
) *UserScheduleDto {
	return &UserScheduleDto{
		userId:       userId,
		fromDateTime: fromDatetime,
		toDateTime:   toDateTime,
		tagIds:       tagIds,
		latitude:     latitude,
		longitude:    longitude,
	}
}

type UsTagsJoinedDto struct {
	UserScheduleId int64           `db:"userScheduleId"`
	UserId         string          `db:"userId"`
	FromDateTime   time.Time       `db:"fromDateTime"`
	ToDateTime     time.Time       `db:"toDateTime"`
	LocationTypeID int8            `db:"locationTypeId"`
	Latitude       sql.NullFloat64 `db:"latitude"`
	Longitude      sql.NullFloat64 `db:"longitude"`
	TagId          sql.NullInt32   `db:"tagId"`
}

// IUserScheduleQueryRepository is an interface for userschedules table
type IUserScheduleQueryRepository interface {
	QueryUserSchedulesWhereTimeRange(beginDateTime, endDateTime time.Time, userId string) ([]*UserScheduleDto, error)
	QueryUserScheduleWhereId(userScheduleId int64) (*UserScheduleDto, error)
}

var _ IUserScheduleQueryRepository = (*realUserScheduleQueryRepository)(nil)

type realUserScheduleQueryRepository struct {
	db SqlDb
}

func ProvideUserScheduleRepository(db SqlDb) IUserScheduleQueryRepository {
	return &realUserScheduleQueryRepository{
		db: db,
	}
}

func (r *realUserScheduleQueryRepository) QueryUserSchedulesWhereTimeRange(beginDateTime, endDateTime time.Time, userId string) ([]*UserScheduleDto, error) {
	var rows *sqlx.Rows
	baseQuery := `
		SELECT us.userScheduleId, us.userId,
			   us.fromDateTime, us.toDateTime,
			   us.locationTypeId,
			   ust.tagId,
			   usl.latitude, usl.longitude
		FROM userschedules us
		LEFT JOIN userschedulelocations usl ON us.userScheduleId=usl.userScheduleId
		LEFT JOIN userscheduletags ust ON us.userScheduleId=ust.userScheduleId
		WHERE fromDateTime >= ? AND toDateTime <= ?
	`
	var (
		err error
	)
	if userId != "" {
		rows, err = r.db.Queryx(baseQuery+" AND userId = ? ORDER BY userId", beginDateTime, endDateTime, userId)
	} else {
		rows, err = r.db.Queryx(baseQuery, beginDateTime, endDateTime)
	}
	if err != nil {
		return nil, stew.Wrap(err)
	}
	var joinedDtos []*UsTagsJoinedDto
	for rows.Next() {
		var jDto UsTagsJoinedDto
		if err = rows.StructScan(&jDto); err != nil {
			return nil, stew.Wrap(err)
		}
		joinedDtos = append(joinedDtos, &jDto)
	}

	return r.compressJoinedDtos(joinedDtos), nil
}

func (r *realUserScheduleQueryRepository) QueryUserScheduleWhereId(userScheduleId int64) (*UserScheduleDto, error) {
	var query = `
		SELECT us.userScheduleId, us.userId,
			   us.fromDateTime, us.toDateTime,
			   us.locationTypeId,
			   ust.tagId,
			   usl.latitude, usl.longitude
		FROM userschedules us
		LEFT JOIN userschedulelocations usl ON us.userScheduleId=usl.userScheduleId
		LEFT JOIN userscheduletags ust ON us.userScheduleId=ust.userScheduleId
		WHERE us.userScheduleId = ?
	`
	rows, err := r.db.Queryx(query, userScheduleId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, stew.Wrap(err)
	}
	var joinedDtos []*UsTagsJoinedDto
	for rows.Next() {
		var jDto UsTagsJoinedDto
		if err = rows.StructScan(&jDto); err != nil {
			return nil, stew.Wrap(err)
		}
		joinedDtos = append(joinedDtos, &jDto)
	}

	uScheduleDtos := r.compressJoinedDtos(joinedDtos)
	return uScheduleDtos[0], nil
}

func (r *realUserScheduleQueryRepository) compressJoinedDtos(joinedDtos []*UsTagsJoinedDto) (uScheduleDtos []*UserScheduleDto) {
	// [Strategy]
	// i. すでに存在しているとき
	// ii. 初登場のカテゴリーが出現したときリストに新しく入れる
	for _, jDto := range joinedDtos {
		// Check if the parent ID is already there in the list
		var doesTheParentIdAlreadyExistInReturnList bool
		for _, cateTag := range uScheduleDtos {
			if cateTag.userScheduleId == jDto.UserScheduleId {
				doesTheParentIdAlreadyExistInReturnList = true
			}
		}
		if doesTheParentIdAlreadyExistInReturnList {
			// i.
			for _, usDto := range uScheduleDtos {
				if usDto.userScheduleId == jDto.UserScheduleId {
					usDto.tagIds = append(usDto.tagIds, uint16(jDto.TagId.Int32))
				}
			}
		} else {
			// ii.
			var tagIds = make([]uint16, 0)
			if jDto.TagId.Valid {
				tagIds = append(tagIds, uint16(jDto.TagId.Int32))
			}
			usDto := newUserScheduleDtoForQuery(
				jDto.UserScheduleId, jDto.UserId,
				jDto.FromDateTime, jDto.ToDateTime,
				tagIds,
				jDto.Latitude.Float64, jDto.Longitude.Float64,
				jDto.LocationTypeID,
			)
			uScheduleDtos = append(uScheduleDtos, usDto)
		}
	}
	return
}

type IUserScheduleCommandRepository interface {
	InsertUserSchedule(dto *UserScheduleDto) (int64, error)
	UpdateUserSchedule(userScheduleId int64, dto *UserScheduleDto) (int64, error)
	DeleteUserSchedule(userScheduleId int64) error
}

var _ IUserScheduleCommandRepository = (*realUserScheduleCommandRepository)(nil)

type realUserScheduleCommandRepository struct {
	db SqlDb
}

func ProvideRealUserScheduleUpdateRepository(db SqlDb) IUserScheduleCommandRepository {
	return &realUserScheduleCommandRepository{
		db: db,
	}
}

func (r *realUserScheduleCommandRepository) InsertUserSchedule(dto *UserScheduleDto) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, stew.Wrap(err)
	}

	// Insert into userschedules table
	res, err := tx.Exec(`
		INSERT INTO userschedules
		(userId, fromDateTime, toDateTime) VALUES (?, ?, ?)`,
		dto.userId, dto.fromDateTime, dto.toDateTime)
	if err != nil {
		tx.Rollback()
		return 0, stew.Wrap(err)
	}
	lastInsertedUserScheduleId, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, stew.Wrap(err)
	}

	// Insert into userschedulelocations table
	_, err = tx.Exec(`
		INSERT INTO userschedulelocations
		(userScheduleId, latitude, longitude) VALUES (?, ?, ?)`,
		lastInsertedUserScheduleId, dto.latitude, dto.longitude)
	if err != nil {
		tx.Rollback()
		return 0, stew.Wrap(err)
	}

	// Insert into userscheduletags table
	for _, tagId := range dto.tagIds {
		_, err = tx.Exec(`
			INSERT INTO userscheduletags
			(userScheduleId, tagId) VALUES (?, ?)`,
			lastInsertedUserScheduleId, tagId)
		if err != nil {
			tx.Rollback()
			return 0, stew.Wrap(err)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, stew.Wrap(err)
	}

	return lastInsertedUserScheduleId, nil
}

func (r *realUserScheduleCommandRepository) UpdateUserSchedule(userScheduleId int64, dto *UserScheduleDto) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, stew.Wrap(err)
	}

	// Delete existing userscheduletags table
	_, err = tx.Exec(`
		DELETE FROM userscheduletags
		WHERE userScheduleId = ?`,
		userScheduleId)
	if err != nil {
		tx.Rollback()
		return 0, stew.Wrap(err)
	}
	// Insert into the updated userscheduletags table
	for _, tagId := range dto.tagIds {
		_, err = tx.Exec(`
			INSERT INTO userscheduletags
			(userScheduleId, tagId) VALUES (?, ?)`,
			userScheduleId, tagId)
		if err != nil {
			tx.Rollback()
			return 0, stew.Wrap(err)
		}
	}

	// Update userschedulelocations table
	_, err = tx.Exec(`
		UPDATE userschedulelocations
		SET latitude = ?, longitude = ? WHERE userScheduleId = ?`,
		dto.latitude, dto.longitude, userScheduleId)
	if err != nil {
		tx.Rollback()
		return 0, stew.Wrap(err)
	}

	// Update userschedules table
	_, err = tx.Exec(`
		UPDATE userschedules
		SET fromDateTime = ?, toDateTime = ? WHERE userScheduleId = ?`,
		dto.fromDateTime, dto.toDateTime, userScheduleId)
	if err != nil {
		tx.Rollback()
		return 0, stew.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return 0, stew.Wrap(err)
	}

	return userScheduleId, nil
}

func (r *realUserScheduleCommandRepository) DeleteUserSchedule(userScheduleId int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return stew.Wrap(err)
	}

	// userschedules
	_, err = tx.Exec(`
		DELETE FROM userschedules
		WHERE userScheduleId = ?`,
		userScheduleId)
	if err != nil {
		tx.Rollback()
		return stew.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return stew.Wrap(err)
	}

	return nil
}
