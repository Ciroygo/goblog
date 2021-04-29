package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
	"unicode/utf8"

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
		view.Render(w, articles, "articles.index")
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
		view.Render(w, article, "articles.show")
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

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)

	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "长度只能在3-40个之间"
	}

	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["title"] = "内容需要大于10哥字"
	}

	return errors
}

func (c *ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		_article.Create()

		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+_article.GetStringID())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误ciroy")
		}
	} else {
		view.Render(w, view.D{
			"Title":  title,
			"Body":   body,
			"Errors": errors,
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

	view.Render(w, view.D{
		"Title":   _article.Title,
		"Body":    _article.Body,
		"Article": _article,
		"Errors":  nil,
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

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errors := validateArticleFormData(title, body)

	if len(errors) == 0 {
		_article.Title = title
		_article.Body = body

		rowsAffected, err := _article.Update()

		if err != nil {
			logger.LogError(err)
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
			"Title":   _article.Title,
			"Body":    _article.Body,
			"Article": _article,
			"Errors":  errors,
		}, "article.edit", "article._form_field")
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
