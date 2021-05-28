package article

import (
	"goblog/app/models"
	"goblog/app/models/user"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"strconv"
)

type Article struct {
	models.BaseModel

	Title string
	Body  string

	UserID uint64 `gorm:"not null;index"`
	User   user.User
}

func (a Article) Link() string {
	return route.Name2URL("articles.show", "id", strconv.FormatInt(int64(a.ID), 10))
}

// GetStringID 获取 ID 的字符串格式
func (a Article) GetStringID() string {
	return types.Int64ToString(int64(a.ID))
}
