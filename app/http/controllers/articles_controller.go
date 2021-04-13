package controllers

import (
	"database/sql"
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
)

// ArticlesController 文章相关
type ArticlesController struct {
}

// Show 文章详情页面
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章没有找到")
		} else {
			logger.LogError(err)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	}

	//tmpl, err := template.ParseFiles("resources/views/articles/show.gohtml")

	tmpl, err := template.New("show.gohtml").Funcs(template.FuncMap{
		"RouteName2URL": route.Name2URL,
		"Int64ToString": types.Int64ToString,
	}).ParseFiles("resources/views/articles/show.gohtml")

	logger.LogError(err)

	tmpl.Execute(w, article)
}
