# golang 教程学习笔记
## tmp 目录是air文件生成

## 第三章
### 设置标头

## 第四章

### 路由 http.ServeMux

### 集成 Gorilla ServeMux

### 依赖管理 Go Modules

### 弃用 $GOPATH

* Go 源码必须放置于 $GOPATH/src 下，抛弃 $GOPATH 的好处，是你能在任意地方创建的 Go 项目。
* 有非常落后的依赖管理系统，无法传达任何版本信息

## 使用中间件

todo 4.4 0305


## 第五章

### 5.5 模版语法

### `{{ }}`
### with 关键字

```golang
{{ with pipeline }} T1 {{ end }}
{{ with pipeline }} T1 {{ else }} T0 {{ end }}
```
