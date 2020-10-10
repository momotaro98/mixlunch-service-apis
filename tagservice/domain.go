package tagservice

import (
	"github.com/momotaro98/stew"
)

type TagType int8

const (
	All      TagType = iota // 0
	Interest                // 1
	Skill                   // 2
)

func (tt TagType) String() string {
	switch tt {
	case All:
		return "All"
	case Interest:
		return "Interest"
	case Skill:
		return "Skill"
	default:
		return ""
	}
}

type CategoryTags struct {
	Category *Category   `json:"category"`
	Tags     []*SmallTag `json:"tags"`
}

func NewCategoryTags(category *Category, tags []*SmallTag) *CategoryTags {
	return &CategoryTags{
		Category: category,
		Tags:     tags,
	}
}

type Category struct {
	CategoryId uint8  `json:"id"`
	Name       string `json:"name"`
}

func NewCategory(categoryId uint8, name string) *Category {
	return &Category{
		CategoryId: categoryId,
		Name:       name,
	}
}

type SmallTag struct {
	TagId uint16 `json:"id"`
	Name  string `json:"name"`
}

func NewSmallTag(tagId uint16, name string) *SmallTag {
	return &SmallTag{
		TagId: tagId,
		Name:  name,
	}
}

type TagServer interface {
	GetTagsByTagType(tagType TagType) ([]*CategoryTags, error)
	GetTagsByTagTypeAndTagIds(tagType TagType, tagIds []uint16) ([]*CategoryTags, error)
}

type RealTagServer struct {
	tagQueryRepository ITagQueryRepository
}

func ProvideTagServer(queryRepository ITagQueryRepository) TagServer {
	return &RealTagServer{
		tagQueryRepository: queryRepository,
	}
}

func (s *RealTagServer) handleFromTagDtosToTagModels(tagDtos []*TagQueryDto) []*CategoryTags {
	var categoryTagsList = make([]*CategoryTags, 0)
	// [Strategy]
	// i. すでに存在しているカテゴリーのときは、CategoryTagsのTagsリストへ"タグだけ入れる"
	// ii. 初登場のカテゴリーが出現したとき[]*CategoryTagsのリストに新しく、"タグと一緒に入れる"
	for _, tagQDto := range tagDtos {
		// Check if the category is already there in the list
		var doesTheCategoryAlreadyExistInTheListSupposedBeReturned bool
		for _, cateTag := range categoryTagsList {
			if cateTag.Category.CategoryId == tagQDto.categoryId {
				doesTheCategoryAlreadyExistInTheListSupposedBeReturned = true
			}
		}
		if doesTheCategoryAlreadyExistInTheListSupposedBeReturned { // [Strategy] i.
			// New Tag
			tag := NewSmallTag(tagQDto.tagId, tagQDto.tagName)
			// Insert into the CategoryTags' one with the new tag
			for _, cateTag := range categoryTagsList {
				if cateTag.Category.CategoryId == tagQDto.categoryId {
					cateTag.Tags = append(cateTag.Tags, tag)
				}
			}
		} else { // [Strategy] ii.
			// New Tag
			tag := NewSmallTag(tagQDto.tagId, tagQDto.tagName)
			// New Category
			category := NewCategory(tagQDto.categoryId, tagQDto.categoryName)
			// New CategoryTags
			categoryTags := NewCategoryTags(category, []*SmallTag{tag})
			// Insert into []*CategoryTags with the new category and the tag
			categoryTagsList = append(categoryTagsList, categoryTags)
		}
	}
	return categoryTagsList
}

// GetTagsByTagType gets tags by using tag type.
// All tags can be got when passed TagType is "All".
func (s *RealTagServer) GetTagsByTagType(tagType TagType) (categoryTagsList []*CategoryTags, err error) {
	// Query tags
	tagQueryDtos, err := s.tagQueryRepository.QueryTagsWhereTagType(uint8(tagType))
	if err != nil {
		return nil, stew.Wrap(err)
	}
	return s.handleFromTagDtosToTagModels(tagQueryDtos), nil
}

// GetTagsByTagTypeAndTagIds gets tags by using tag type and tag IDs.
// All tags can be got when passed TagType is "All".
func (s *RealTagServer) GetTagsByTagTypeAndTagIds(tagType TagType, tagIds []uint16) ([]*CategoryTags, error) {
	// [Issue]
	// For now the method gets all of tag rows (extend to in-memory)
	// then filters by specified tag IDs.
	// Considering performance (large in-memory allocation), it should leave the filtering task to Database.

	// Query tags
	tagQueryDtos, err := s.tagQueryRepository.QueryTagsWhereTagType(uint8(tagType))
	if err != nil {
		return nil, stew.Wrap(err)
	}
	// Do filtering by specified tag IDs
	var filteredTagQueryDtos = make([]*TagQueryDto, 0)
	for _, tqDtos := range tagQueryDtos {
		for _, tagId := range tagIds {
			if tagId == tqDtos.tagId {
				filteredTagQueryDtos = append(filteredTagQueryDtos, tqDtos)
			}
		}
	}
	return s.handleFromTagDtosToTagModels(filteredTagQueryDtos), nil
}
