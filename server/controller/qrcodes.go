package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func QrcodesController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256) //把文本变成二进制图片
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "image/png", png) //把图片传给前端（不下载，展示）
	} else {
		c.Status(http.StatusBadRequest)
	}
}
