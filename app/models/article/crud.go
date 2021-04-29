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

	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

func GetAll() ([]Article, error) {
	var articles []Article

	if err := model.DB.Find(&articles).Error; err != nil {
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
