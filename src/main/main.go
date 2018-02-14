package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var subProcess *exec.Cmd

	//退出时杀掉启动的进程
	defer func() {
		// 判断子进程是否已经自动退出
		if subProcess.ProcessState.Exited() {
			fmt.Println("")
			fmt.Println("~~~项目已经退出~~~")
			return
		}

		// 不用手动杀死进程，ctrl+c 终止主进程时同时会终止子进程

		// 杀掉子进程
		// err := subProcess.Process.Kill()
		// if err != nil {
		// 	fmt.Println("项目进程退出失败:" + err.Error())
		// } else {
		// 	fmt.Println("项目进程退出成功")
		// }

		// 杀掉所有同名的进程
		// cmdKillAll := exec.Command("killall", binFile)
		// err := cmdKillAll.Run()
		// if err != nil {
		// 	fmt.Println("项目进程退出失败:" + err.Error())
		// } else {
		// 	fmt.Println("项目进程退出成功")
		// }
		return
	}()

	//获得当前路径
	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 寻找根路径
	rootPath := FindRoot(curPath)

	//判断是否在项目的根目录
	findPath := rootPath + "/src/main/main.go"
	_, err = os.Stat(findPath)
	if err != nil {
		fmt.Println(findPath+"不存在", err.Error())
		return
	}

	//将当前路径设置为GOPATH
	err = os.Setenv("GOPATH", rootPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("GOPATH设置完成")

	//获得当前路径的项目名
	rootPaths := strings.Split(rootPath, "/")
	binFile := rootPaths[len(rootPaths)-1] + "_.gorun"

	fmt.Println("编译项目...")
	cmd := exec.Command("go", "build", "-o", binFile, "./src/main/")
	cmd.Dir = rootPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("编译项目完成")
	fmt.Println("~~~开始运行项目~~~ 以下开始为项目内的输出: ")
	fmt.Println("")
	subProcess = exec.Command("./" + binFile)
	subProcess.Dir = rootPath
	subProcess.Stderr = os.Stderr
	subProcess.Stdout = os.Stdout
	err = subProcess.Run()
	if err != nil {
		fmt.Println("运行项目失败: " + err.Error())
		return
	}
}

// FindRoot 用来找到项目根路径
func FindRoot(path string) string {
	var mark = false
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("终端所在路径错误", err.Error())
		return ""
	}

	for _, v := range files {
		if v.IsDir() && v.Name() == "src" {
			mark = true
			break
		}
	}

	if !mark {
		index := strings.LastIndex(path, "/")
		path := path[:index]
		return FindRoot(path)
	}
	return path
}
