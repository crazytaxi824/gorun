package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
)

var subProcess *exec.Cmd

func main() {
	go KillProcess(subProcess)

	//获得当前路径
	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("%c[0;0;31m%s%c[0m", 0x1B, "当前路径获取失败: ", 0x1B)
		fmt.Println(err.Error())
		return
	}

	// 向上层文件夹寻找根路径
	rootPath, err := FindRoot(curPath)
	if err != nil {
		fmt.Printf("%c[0;0;31m%s%c[0m", 0x1B, err.Error(), 0x1B)
		fmt.Println("")
		return
	}

	//将当前路径设置为GOPATH
	err = os.Setenv("GOPATH", rootPath)
	if err != nil {
		fmt.Printf("%c[0;0;31m%s%c[0m", 0x1B, "设置GOPATH错误: ", 0x1B)
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%c[0;0;32m%s%c[0m", 0x1B, "GOPATH设置完成", 0x1B)
	fmt.Println("")

	//获得当前路径的项目名
	rootPaths := strings.Split(rootPath, "/")
	gorunFile := rootPaths[len(rootPaths)-1] + "_.gorun"

	fmt.Printf("%c[0;0;32m%s%c[0m", 0x1B, "编译项目...", 0x1B)
	fmt.Println("")
	cmd := exec.Command("go", "build", "-o", gorunFile, "./src/main/")
	cmd.Dir = rootPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%c[0;0;31m%s%c[0m", 0x1B, "编译项目错误: ", 0x1B)
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%c[0;0;32m%s%c[0m", 0x1B, "编译项目完成", 0x1B)
	fmt.Println("")
	fmt.Printf("%c[0;0;32m%s%c[0m", 0x1B, "~~~开始运行项目~~~ 以下开始为项目内的输出: ", 0x1B)
	fmt.Println("")
	fmt.Println("")
	subProcess = exec.Command("./" + gorunFile)
	subProcess.Dir = rootPath
	subProcess.Stderr = os.Stderr
	subProcess.Stdout = os.Stdout
	err = subProcess.Start()
	if err != nil {
		fmt.Printf("%c[0;0;31m%s%c[0m", 0x1B, "运行项目失败: ", 0x1B)
		fmt.Println(err.Error())
		return
	}
	err = subProcess.Wait()
	if err != nil {
		fmt.Printf("%c[0;0;31m%s%c[0m", 0x1B, "请检查 "+gorunFile+" : "+err.Error(), 0x1B)
		fmt.Println("")
		fmt.Println("")
		return
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Printf("%c[0;0;32m%s%c[0m", 0x1B, "~~~项目已经退出~~~", 0x1B)
	fmt.Println("")
	fmt.Println("")
}

// FindRoot 用来找到项目根路径
func FindRoot(path string) (string, error) {
	var mark = false
	_, err := os.Stat(path + "/src/main/main.go")
	if err == nil {
		mark = true
	}

	if !mark {
		index := strings.LastIndex(path, "/")
		if index == -1 {
			return "", errors.New("执行路径错误")
		}
		path := path[:index]
		return FindRoot(path)
	}
	return path, nil
}

// KillProcess 杀进程
func KillProcess(subProcess *exec.Cmd) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
	fmt.Println("")
	fmt.Printf("%c[0;0;31m%s%c[0m", 0x1B, "程序被强制中断!!!", 0x1B)
	fmt.Println("")
}
