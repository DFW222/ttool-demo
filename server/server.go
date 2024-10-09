package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
	c "ttool/server/controller" //给包取别名为c
	"ttool/server/ws"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run() {
	hub := ws.NewHub()
	go hub.Run()

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	staticFiles, _ := fs.Sub(FS, "frontend/dist") //变成结构化目录（把静态文件变成一个变量）
	router.POST("/api/v1/files", c.FilesController)
	router.GET("/api/v1/qrcodes", c.QrcodesController)
	router.GET("/uploads/:path", c.UploadsController)
	router.GET("/api/v1/addresses", c.AddressesController)
	router.POST("/api/v1/texts", c.TextsController)
	router.GET("/ws", func(c *gin.Context) {
		ws.HttpController(c, hub)
	})
	router.StaticFS("/static", http.FS(staticFiles)) //第一个参数为静态文件的前缀，所有静态文件都放在这个前缀后面，用http读取dist内的文件
	router.NoRoute(func(c *gin.Context) {            //文件处理
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	router.Run(":27149")
}
