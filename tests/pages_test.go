package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHomePage(t *testing.T)  {
	bashURL := "http://localhost:3000"

	var (
		resp *http.Response
		err error
	)

	resp, err = http.Get(bashURL + "/")

	assert.NoError(t, err, "有错误发生 err不为空")
	assert.Equal(t, 200, resp.StatusCode, "应该返回状态码 200")
}