package tagservice

import (
	"database/sql"

	"github.com/momotaro98/stew"
)

type TagQueryDto struct {
	tagId        uint16
	tagName      string
	categoryId   uint8
	categoryName string
}

type ITagQueryRepository interface {
	QueryTagsWhereTagType(tagTypeId uint8) ([]*TagQueryDto, error)
}

var _ ITagQueryRepository = (*realTagQueryRepository)(nil)

type realTagQueryRepository struct {
	db SqlDb
}

func ProvideTagQueryRepository(db SqlDb) ITagQueryRepository {
	return &realTagQueryRepository{
		db: db,
	}
}

// QueryTagsWhereTagType does query Tags master.
// All tags (tag types) can be got when the passed tagTypeId is 0.
func (r *realTagQueryRepository) QueryTagsWhereTagType(tagTypeId uint8) ([]*TagQueryDto, error) {
	var err error
	var rows *sql.Rows
	if tagTypeId == 0 {
		rows, err = r.db.Query(`
		SELECT t.tagId, t.name, c.categoryId, c.name
		FROM tags t
		INNER JOIN categories c ON t.categoryId = c.categoryId
		ORDER BY c.categoryId ASC`)
	} else {
		rows, err = r.db.Query(`
		SELECT t.tagId, t.name, c.categoryId, c.name
		FROM tags t
		INNER JOIN categories c ON t.categoryId = c.categoryId
		WHERE t.tagTypeId = ?
		ORDER BY c.categoryId ASC`,
			tagTypeId)
	}
	if err != nil {
		return nil, stew.Wrap(err)
	}
	var tagDtos []*TagQueryDto
	for rows.Next() {
		var tDto TagQueryDto
		if err := rows.Scan(&tDto.tagId, &tDto.tagName, &tDto.categoryId, &tDto.categoryName); err != nil {
			return nil, stew.Wrap(err)
		}
		tagDtos = append(tagDtos, &tDto)
	}
	return tagDtos, nil
}
