package partyservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/momotaro98/stew"
)

type PartyDto struct {
	id         int64
	startFrom  time.Time
	endTo      time.Time
	chatRoomId sql.NullString
}

type PartyMemberDto struct {
	partyMemberId int64
	partyId       int64
	userId        string
}

type PartyTagsDto struct {
	partyId int64
	tagIds  []uint16
}

type IPartyQueryRepository interface {
	QueryPartiesWhereTimeRange(queryDto *PartyQueryDto) ([]*PartyDto, error)
	QueryPartiesWhereUserIdAndTimeRange(userId string, queryDto *PartyQueryDto) ([]*PartyDto, error)
	QueryPartiesWhereUserIdLastN(userId string, n int) ([]*PartyDto, error)
	QueryPartyMembersWherePartyId(partyId int64) ([]*PartyMemberDto, error)
	QueryPartyTagsWherePartyId(partyId int64) (*PartyTagsDto, error)
	QueryPartyReviewMembers(queryDto *ReviewMemberQueryDto) ([]*PartyMemberReviewDto, error)
}

var _ IPartyQueryRepository = (*realPartyQueryRepository)(nil)

type realPartyQueryRepository struct {
	db SqlDb
}

func ProvidePartyQueryRepository(db SqlDb) IPartyQueryRepository {
	return &realPartyQueryRepository{
		db: db,
	}
}

type PartyQueryDto struct {
	beginDateTime *time.Time
	endDateTime   *time.Time
}

func buildSQLForQueryPartiesWhereTimeRange(queryDto *PartyQueryDto) (sql string, args []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(
		"id", "startFrom", "endTo", "chatRoomId",
	)
	sb.From("parties")
	if dt := queryDto.beginDateTime; dt != nil {
		sb.Where(sb.GreaterEqualThan("startFrom", *dt))
	}
	if dt := queryDto.endDateTime; dt != nil {
		sb.Where(sb.LessEqualThan("endTo", *dt))
	}
	return sb.Build()
}

func (r *realPartyQueryRepository) QueryPartiesWhereTimeRange(queryDto *PartyQueryDto) ([]*PartyDto, error) {
	query, args := buildSQLForQueryPartiesWhereTimeRange(queryDto)
	return r.queryPartyDtos(query, args...)
}

func buildSQLForQueryPartiesWhereUserIdAndTimeRange(userId string, queryDto *PartyQueryDto) (sql string, args []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(
		sb.As("p.id", "id"), sb.As("p.startFrom", "startFrom"),
		sb.As("p.endTo", "endTo"), sb.As("p.chatRoomId", "chatRoomId"),
	)
	sb.From(sb.As("partymembers", "pm"))
	sb.Join(sb.As("parties", "p"), "pm.partyId = p.id")
	sb.Where(sb.Equal("pm.userId", userId))
	if dt := queryDto.beginDateTime; dt != nil {
		sb.Where(sb.GreaterEqualThan("p.startFrom", *dt))
	}
	if dt := queryDto.endDateTime; dt != nil {
		sb.Where(sb.LessEqualThan("p.endTo", *dt))
	}
	return sb.Build()
}

func (r *realPartyQueryRepository) QueryPartiesWhereUserIdAndTimeRange(userId string, queryDto *PartyQueryDto) ([]*PartyDto, error) {
	query, args := buildSQLForQueryPartiesWhereUserIdAndTimeRange(userId, queryDto)
	return r.queryPartyDtos(query, args...)
}

func (r *realPartyQueryRepository) QueryPartiesWhereUserIdLastN(userId string, n int) ([]*PartyDto, error) {
	return r.queryPartyDtos(`
		SELECT p.id, p.startFrom, p.endTo, p.chatRoomId
		FROM partymembers pm
		INNER JOIN parties p ON pm.partyId = p.id
		WHERE pm.userId = ?
		ORDER BY p.id DESC
		LIMIT ?
`,
		userId, n,
	)
}

func (r *realPartyQueryRepository) queryPartyDtos(query string, args ...interface{}) ([]*PartyDto, error) {
	var partyDtos []*PartyDto
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	for rows.Next() {
		var pDto PartyDto
		if err := rows.Scan(&pDto.id, &pDto.startFrom, &pDto.endTo, &pDto.chatRoomId); err != nil {
			return nil, err
		}
		partyDtos = append(partyDtos, &pDto)
	}
	return partyDtos, nil
}

func (r *realPartyQueryRepository) QueryPartyMembersWherePartyId(partyId int64) ([]*PartyMemberDto, error) {
	var pMemberDtos []*PartyMemberDto
	rows, err := r.db.Query(`
		SELECT pm.partyMemberId, pm.userId, pm.partyId
		FROM partymembers pm
		WHERE pm.partyId = ?`,
		partyId)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	for rows.Next() {
		var pmDto PartyMemberDto
		if err := rows.Scan(&pmDto.partyMemberId, &pmDto.userId, &pmDto.partyId); err != nil {
			return nil, err
		}
		pMemberDtos = append(pMemberDtos, &pmDto)
	}
	return pMemberDtos, nil
}

func (r *realPartyQueryRepository) QueryPartyTagsWherePartyId(partyId int64) (*PartyTagsDto, error) {
	rows, err := r.db.Query(`
		SELECT partyId, tagId
		FROM partytags
		where partyId = ?`,
		partyId)
	if err != nil {
		return nil, stew.Wrap(err)
	}

	type ptDto struct {
		partyID int64
		tagID   uint16
	}

	var dtos []*ptDto
	for rows.Next() {
		var dto ptDto
		if err = rows.Scan(
			&dto.partyID,
			&dto.tagID,
		); err != nil {
			return nil, stew.Wrap(err)
		}
		dtos = append(dtos, &dto)
	}

	var partyTagsDto PartyTagsDto
	for _, dto := range dtos {
		partyTagsDto.partyId = dto.partyID
		partyTagsDto.tagIds = append(partyTagsDto.tagIds, dto.tagID)
	}

	return &partyTagsDto, nil
}

type ReviewMemberQueryDto struct {
	partyID  int64
	reviewer string
	reviewee string
}

func buildSQLForQueryPartyReviewMembers(queryDto *ReviewMemberQueryDto) (sql string, args []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("partyId", "reviewer", "reviewee", "score", "comments")
	sb.From("partymemberreviews")
	if queryDto.partyID >= 1 {
		sb.Where(sb.Equal("partyId", queryDto.partyID))
	}
	if queryDto.reviewer != "" {
		sb.Where(sb.Equal("reviewer", queryDto.reviewer))
	}
	if queryDto.reviewee != "" {
		sb.Where(sb.Equal("reviewee", queryDto.reviewee))
	}
	return sb.Build()
}

func (r *realPartyQueryRepository) QueryPartyReviewMembers(queryDto *ReviewMemberQueryDto) ([]*PartyMemberReviewDto, error) {
	// Build query
	query, args := buildSQLForQueryPartyReviewMembers(queryDto)
	// Query
	var retDtos []*PartyMemberReviewDto
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, stew.Wrap(err)
	}
	for rows.Next() {
		var dto PartyMemberReviewDto
		if err := rows.Scan(&dto.partyID, &dto.reviewer, &dto.reviewee, &dto.score, &dto.comments); err != nil {
			return nil, stew.Wrap(err)
		}
		retDtos = append(retDtos, &dto)
	}
	return retDtos, nil
}

type IPartyCommandRepository interface {
	Tran() (*sql.Tx, error)
	Commit(*sql.Tx) error
	Rollback(*sql.Tx) error
	InsertParty(tx *sql.Tx, dto *PartyCommandDto) (int64, error)
	DeletePartiesWithADay(tx *sql.Tx, targetDay time.Time) error
	InsertPartyMemberReview(tx *sql.Tx, dto *PartyMemberReviewDto) error
}

var _ IPartyCommandRepository = (*realPartyCommandRepository)(nil)

type realPartyCommandRepository struct {
	db SqlDb
}

func ProvidePartyCommandRepository(db SqlDb) IPartyCommandRepository {
	return &realPartyCommandRepository{
		db: db,
	}
}

type PartyCommandDto struct {
	id            int64
	startFrom     time.Time
	endTo         time.Time
	chatRoomId    sql.NullString
	memberUserIDs []string
}

func (r *realPartyCommandRepository) Tran() (*sql.Tx, error) {
	return r.db.DB.Begin()
}

func (r *realPartyCommandRepository) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (r *realPartyCommandRepository) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (r *realPartyCommandRepository) InsertParty(tx *sql.Tx, dto *PartyCommandDto) (int64, error) {
	// Inserting to parties table
	res, err := tx.Exec("INSERT INTO parties (startFrom, endTo, chatRoomId) VALUES (?, ?, ?)",
		dto.startFrom, dto.endTo, dto.chatRoomId)
	if err != nil {
		return 0, stew.Wrap(err)
	}
	insertedPartyId, err := res.LastInsertId()
	if err != nil {
		return 0, stew.Wrap(err)
	}

	// Inserting to partymembers table
	for _, userID := range dto.memberUserIDs {
		_, err := tx.Exec("INSERT INTO partymembers (partyId, userId) VALUES (?, ?)",
			insertedPartyId, userID)
		if err != nil {
			return 0, stew.Wrap(err)
		}
	}

	return insertedPartyId, nil
}

func (r *realPartyCommandRepository) DeletePartiesWithADay(tx *sql.Tx, targetDay time.Time) error {
	start := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(),
		0, 0, 0, 0, time.Local)
	end := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(),
		23, 59, 59, 999, time.Local)
	if _, err := tx.Exec(`
		DELETE FROM parties
		WHERE startFrom >= ? AND endTo <= ?`,
		start, end); err != nil {
		return stew.Wrap(err)
	}
	return nil
}

type PartyMemberReviewDto struct {
	partyID  int64
	reviewer string
	reviewee string
	score    float64
	comments string
}

func (r *realPartyCommandRepository) InsertPartyMemberReview(tx *sql.Tx, dto *PartyMemberReviewDto) error {
	_, err := tx.Exec(`
		INSERT INTO partymemberreviews (partyId, reviewer, reviewee, score, comments)
		VALUES (?, ?, ?, ?, ?)`,
		dto.partyID, dto.reviewer, dto.reviewee, dto.score, dto.comments)
	if err != nil {
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
	return nil
}

type IChatRoomRepository interface {
	CreateChatRoom(chatRoomId string) error
}

var _ IChatRoomRepository = (*realChatRoomRepository)(nil)

type realChatRoomRepository struct {
	app App
}

func ProvideChatRoomRepository(app App) IChatRoomRepository {
	return &realChatRoomRepository{
		app: app,
	}
}

const (
	document   = "rooms"
	keyOfChats = "messages"
)

func (r *realChatRoomRepository) CreateChatRoom(chatRoomId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Access to Firebase Cloud Firestore
	// https://firebase.google.com/docs/firestore/manage-data/add-data
	client, err := r.app.Firestore(ctx)
	if err != nil {
		return stew.Wrap(err)
	}
	defer client.Close()

	_, err = client.Collection(document).Doc(chatRoomId).Set(ctx, map[string]interface{}{
		keyOfChats: []string{},
	})
	if err != nil {
		return stew.Wrap(err)
	}
	return nil
}
