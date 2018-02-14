package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("进程: ", os.Getpid(), " 已经启动...")

	//获得当前路径
	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//判断是否在项目的根目录
	findPath := curPath + "/src/main/main.go"
	_, err = os.Stat(findPath)
	if err != nil {
		fmt.Println(findPath+"不存在", err.Error())
		return
	}

	//将当前路径设置为GOPATH
	err = os.Setenv("GOPATH", curPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("GOPATH设置完成")

	//获得当前路径的项目名
	curPaths := strings.Split(curPath, "/")
	binFile := curPaths[len(curPaths)-1] + "_bin"

	fmt.Println("编译项目...")
	cmd := exec.Command("go", "build", "-o", binFile, "./src/main/")
	cmd.Dir = curPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//退出时杀掉启动的进程
	defer func() {
		// 判断子程序是否已经退出
		if cmd.ProcessState.Exited() {
			fmt.Println("子进程", cmd.Process.Pid, "已经退出")
			return
		}

		// 杀掉pid子进程
		i := cmd.Process.Pid
		pidstr := strconv.Itoa(i)
		cmdKill := exec.Command("kill", pidstr)
		err := cmdKill.Run()
		if err != nil {
			fmt.Println("项目进程退出失败:" + err.Error())
		} else {
			fmt.Println("项目进程退出成功")
		}

		// 杀掉所有名字的进程
		// cmdKillAll := exec.Command("killall", binFile)
		// err := cmdKillAll.Run()
		// if err != nil {
		// 	fmt.Println("项目进程退出失败:" + err.Error())
		// } else {
		// 	fmt.Println("项目进程退出成功")
		// }
	}()

	fmt.Println("编译项目完成，开始运行项目...")
	fmt.Println("以下开始为项目内的输出：")
	cmd = exec.Command("./" + binFile)
	cmd.Dir = curPath
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println("运行项目失败：" + err.Error())
		return
	}
	fmt.Println("子进程 pid: ", cmd.Process.Pid, " 已经启动...")
}
