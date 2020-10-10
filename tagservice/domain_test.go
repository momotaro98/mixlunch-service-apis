package tagservice

import (
	"testing"

	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source=repositories.go -destination=repositories_mock.go -package=tagservice -self_package=github.com/momotaro98/mixlunch-service-api/tagservice

func TestGetTagsByTagType_RegularAllCase_Success(t *testing.T) {
	// Arrange
	// DTO to return in mock
	tqDto1 := &TagQueryDto{
		tagId:        1,
		tagName:      "Vue.js",
		categoryId:   1,
		categoryName: "Programming",
	}
	tqDto2 := &TagQueryDto{
		tagId:        2,
		tagName:      "Python",
		categoryId:   1,
		categoryName: "Programming",
	}
	tqDto3 := &TagQueryDto{
		tagId:        3,
		tagName:      "Fishing",
		categoryId:   2,
		categoryName: "Hobby",
	}
	// mock
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	tagQueryRepositoryMock := NewMockITagQueryRepository(mockCtrl)
	tagQueryRepositoryMock.EXPECT().
		QueryTagsWhereTagType(uint8(All)).
		Return([]*TagQueryDto{tqDto1, tqDto2, tqDto3}, nil)
	tagServer := ProvideTagServer(tagQueryRepositoryMock)
	// Act
	categoriesTags, _ := tagServer.GetTagsByTagType(All)
	// Assert
	if len(categoriesTags) != 2 {
		t.Errorf("2 dayo")
	}
}
