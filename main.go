package main

import (
	"os"
	"os/exec"
	"os/signal"
	"ttool/server"
)

func main() {
	chDie := make(chan struct{})
	go server.Run()            //启动gin协程，初始化gin框架
	go startBrowser(chDie)     //打开浏览器
	chs := listenToInterrupt() //监听中断信号
	select {
	case <-chs:
	case <-chDie:
		os.Exit(0) //没信息就一直等，有信息就从chs中发送出去，阻塞结束，执行杀死进程
	}
}

func startBrowser(chDie chan struct{}) {
	//通过执行命令的os/exec库打开chrome浏览器
	chromePath := "/usr/bin/google-chrome"
	cmd := exec.Command(chromePath, "--app=http://localhost:27149/static/index.html")
	cmd.Start()
	cmd.Wait()
	chDie <- struct{}{}
}

func listenToInterrupt() chan os.Signal {
	//创建一个channel,chs变量可以接收一个系统信号用于中断
	chs := make(chan os.Signal, 1)
	signal.Notify(chs, os.Interrupt) //如果触发中断，往chs中发送信息
	return chs
}
