package article

import (
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToInt(idstr)

	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

func GetAll() ([]Article, error) {
	var articles []Article

	if err := model.DB.Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

func (article Article) Update() (rowsAffected int64, err error) {
	result := model.DB.Save(&article)

	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil
}

func (a *Article) Create() (err error) {

	if err = model.DB.Create(&a).Error; err != nil {
		logger.LogError(err)
		fmt.Println(err)
		return err
	}

	return nil
}

func (article *Article) Delete() (RowsAffected int64, err error) {
	result := model.DB.Delete(&article)

	if err = result.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return result.RowsAffected, nil
}

// CreatedAtDate 创建日期
func (a *Article) CreatedAtDate() string {
	return a.CreatedAt.Format("2006-01-02")
}

// GetByUserID 获取全部文章
func GetByUserID(uid string) ([]Article, error) {
	var articles []Article
	if err := model.DB.Where("user_id = ?", uid).Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}
