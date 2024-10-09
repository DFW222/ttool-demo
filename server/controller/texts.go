package controller

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TextsController(c *gin.Context) {
	var json struct { //容纳用户上传的文本接口
		Raw string
	}
	if err := c.ShouldBindJSON(&json); err != nil { //主要绑定成功，就把text请求发给我
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() //获取当前可执行文件目录
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe) //获取当前执行文件目录
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()          //生成一个文件名
		uploads := filepath.Join(dir, "uploads") //拼接uploads的绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  //创建uploads目录，权限等级
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")                        //拼接文件的绝对路径（不含目录）
		err = os.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) //将json.Raw写入文件
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath}) //返回文件的绝对路径（不含exe目录）
	}

}
