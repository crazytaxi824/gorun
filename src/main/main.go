package main

import (
	"errors"
	"fmt"
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
			os.Exit(0)
		}

		// 不用手动杀死进程，ctrl+c 终止主进程时同时会终止子进程

		// 杀掉单独的子进程
		// err := subProcess.Process.Kill()
		// if err != nil {
		// 	fmt.Println("项目进程退出失败:" + err.Error())
		// } else {
		// 	fmt.Println("项目进程退出成功")
		// }

		// 杀掉所有同名的所有进程
		// cmdKillAll := exec.Command("killall", binFile)
		// err := cmdKillAll.Run()
		// if err != nil {
		// 	fmt.Println("项目进程退出失败:" + err.Error())
		// } else {
		// 	fmt.Println("项目进程退出成功")
		// }
		os.Exit(0)
	}()

	//获得当前路径
	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// 向上层文件夹寻找根路径
	rootPath, err := FindRoot(curPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//将当前路径设置为GOPATH
	err = os.Setenv("GOPATH", rootPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("GOPATH设置完成")

	//获得当前路径的项目名
	rootPaths := strings.Split(rootPath, "/")
	gorunFile := rootPaths[len(rootPaths)-1] + "_.gorun"

	fmt.Println("编译项目...")
	cmd := exec.Command("go", "build", "-o", gorunFile, "./src/main/")
	cmd.Dir = rootPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("编译项目完成")
	fmt.Println("~~~开始运行项目~~~ 以下开始为项目内的输出: ")
	fmt.Println("")
	subProcess = exec.Command("./" + gorunFile)
	subProcess.Dir = rootPath
	subProcess.Stderr = os.Stderr
	subProcess.Stdout = os.Stdout
	err = subProcess.Run()
	if err != nil {
		fmt.Println("运行项目失败: " + err.Error())
		os.Exit(1)
	}
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
