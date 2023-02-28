package main

import (
	"fmt"
	"gee/base"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type student struct {
	Name string
	Age int
}

func FormatAsDate(t time.Time) string {
	y,m,d := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
}

func only4V2() base.HandlerFunc {
	return func(c *base.Context) {
		t := time.Now()
		c.Fail(500, "服务器错误")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	e := base.New()
	e.Use(base.Logger())

	e.SetFuncMap(template.FuncMap{"FormatAsDate": FormatAsDate})
	e.LoadHTMLGlob("templates/*")
	e.Static("/assets", "./static")

	stu1 := &student{"lucy", 19}
	stu2 := &student{"lily", 20}

	e.Get("/index", func(c *base.Context) {
		c.Html(http.StatusOK, "css.tmpl", nil)
	})
	e.Get("/students", func(c *base.Context) {
		c.Html(http.StatusOK, "arr.tmpl", base.H{
			"title": "gee",
			"data": [2]*student{stu1, stu2},
		})
	})
	e.Get("/date", func(c *base.Context) {
		c.Html(http.StatusOK, "custom.tmpl", base.H{
			"title": "gee",
			"now": time.Date(2019,8,17,0,0,0,0, time.UTC),
		})
	})
	e.Get("/panic", func(c *base.Context) {
		names := make([]string, 3, 10)
		names[0] = "hello"
		names[1] = "world"
		names[2] = "shit"
		c.String(http.StatusOK, strings.Join(names[:5]," "))
		//如果切片操作超出cap(s)的上限将导致一个panic异常，但是超出len(s)则是意味着扩展了slice，所以会有如下错误，上面代码则不会报错
		//2022/09/20 15:55:15 http: panic serving 127.0.0.1:59349: runtime error: index out of range [11] with length 3
		//c.String(http.StatusOK, names[11])
	})

	v1 := e.Group("/v1")
	//v1.Get("/", func(c *base.Context) {
	//	c.Html(http.StatusOK, "<h1>This is homepage</h1>")
	//})
	v1.Get("/hello", func(c *base.Context) {
		c.String(http.StatusOK, "%s Welcome to go world", c.Query("name"))
	})

	v2 := e.Group("/v2")
	v2.Use(only4V2())

	v2.Get("/hello/:name", func(c *base.Context) {
		c.String(http.StatusOK, "%s Welcome to go world, you are at %s", c.Param("name"), c.Path)
	})
	v2.Get("/assets/*filepath", func(c *base.Context) {
		c.Json(http.StatusOK, base.H{
			"filepath": c.Param("filepath"),
		})
	})
	e.Post("/login", func(c *base.Context) {
		c.Json(http.StatusOK, base.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	e.Run("localhost:9999")
}