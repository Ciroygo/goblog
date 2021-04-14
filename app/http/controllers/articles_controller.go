package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"unicode/utf8"
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
		viewDir := "resources/views"

		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")

		newFiles := append(files, viewDir+"/articles/index.gohtml")

		tmpl, err := template.ParseFiles(newFiles...)

		logger.LogError(err)

		tmpl.ExecuteTemplate(w, "app", articles)
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

		viewDir := "resources/views"

		files, err := filepath.Glob(viewDir + "/layouts/*.gohtml")
		logger.LogError(err)

		newFiles := append(files, viewDir+"/articles/show.gohtml")

		tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
			"Int64ToString": types.Int64ToString,
		}).ParseFiles(newFiles...)

		logger.LogError(err)

		tmpl.ExecuteTemplate(w, "app", article)
	}
}

type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

// 文章保存
func (c *ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")

	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, data)
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
			fmt.Fprint(w, "插入成功，ID为"+strconv.FormatInt(int64(_article.ID), 10))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误ciroy")
		}
	} else {
		storeURL := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
		if err != nil {
			panic(err)
		}

		tmpl.Execute(w, data)
	}
}

func (c *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	article, err := article.Get(id)

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

	updateURL := route.Name2URL("articles.update", "id", id)

	data := ArticlesFormData{
		Title:  article.Title,
		Body:   article.Body,
		URL:    updateURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
	logger.LogError(err)

	tmpl.Execute(w, data)

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
		updateURL := route.Name2URL("articles.update", "id", id)

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    updateURL,
			Errors: errors,
		}

		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)

		tmpl.Execute(w, data)
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
