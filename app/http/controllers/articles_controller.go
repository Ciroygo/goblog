package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"

	"gorm.io/gorm"
)

// ArticlesController 文章相关
type ArticlesController struct {
}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		view.Render(w, view.D{
			"Articles": articles,
		}, "articles.index", "articles._article_meta")
	}
}

// Show 文章详情页面
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没有找到")
		} else {
			logger.LogError(err)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		view.Render(w, view.D{
			"Article":          article,
			"CanModifyArticle": policies.CanModifyArticle(article),
		}, "articles.show", "articles._article_meta")
	}
}

// type ArticlesFormData struct {
// 	Title, Body string
// 	URL         string
// 	Article     article.Article
// 	Errors      map[string]string
// }

// 文章保存
func (c *ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "articles.create")
}

func (c *ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	_article := article.Article{
		Title: r.PostFormValue("title"),
		Body:  r.PostFormValue("body"),
	}
	errors := requests.ValidateArticleForm(_article)

	if len(errors) == 0 {
		_article.Create()

		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误ciroy")
		}
	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errors,
		}, "articles.create")
	}
}

func (c *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 ciroy 文章没有找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器粗错")
		}
		return
	}

	if !policies.CanModifyArticle(_article) {
		flash.Warning("未授权操作！")
		http.Redirect(w, r, "/", http.StatusFound)
	}

	view.Render(w, view.D{
		"Article": _article,
		"Errors":  view.D{},
	}, "articles.edit", "articles._form_field")
}

func (c *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	}
	if !policies.CanModifyArticle(_article) {
		flash.Warning("未授权操作！")
		http.Redirect(w, r, "/", http.StatusForbidden)
	}

	// 4.1 表单验证
	_article.Title = r.PostFormValue("title")
	_article.Body = r.PostFormValue("body")

	errors := requests.ValidateArticleForm(_article)

	if len(errors) == 0 {
		rowsAffected, err := _article.Update()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
			return
		}

		if rowsAffected > 0 {
			//showURL := route.Name2URL("articles.show", "id", id)
			//http.Redirect(w, r, showURL, http.StatusFound)
			showURL := route.Name2URL("articles.show", "id", id)
			http.Redirect(w, r, showURL, http.StatusFound)
		} else {
			fmt.Fprint(w, "您没有做任何修改")
		}

	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errors,
		}, "articles.edit", "articles._form_field")
	}
}

func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章没有找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}

		return
	}

	// 检查权限
	if !policies.CanModifyArticle(article) {
		flash.Warning("您没有权限执行此操作！")
		http.Redirect(w, r, "/", http.StatusFound)
	}

	rowAffected, err := article.Delete()

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 删除遇到故障")
		return
	}

	if rowAffected > 0 {
		indexURL := route.Name2URL("articles.index")
		http.Redirect(w, r, indexURL, http.StatusFound)
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 文章没哟找到")
}
