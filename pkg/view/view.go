package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

func Render(w io.Writer, data interface{}, tplFiles ...string) {
	// 1 设置模板相对路径
	viewDir := "resources/views/"

	// 2. 语法糖，将 articles.show 更正为 articles/show
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".gohtml"
	}

	// 3 所有布局模板文件 Slice
	layoutsFiles, err := filepath.Glob(viewDir + "layouts/*.gohtml")
	logger.LogError(err)

	// 4 在 Slice 里新增我们的目标文件
	allFiles := append(layoutsFiles, tplFiles...)

	// 5 解析所有模版文锦啊

	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)

	logger.LogError(err)

	tmpl.ExecuteTemplate(w, "app", data)
}
