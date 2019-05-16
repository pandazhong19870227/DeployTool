package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

var (
	Debug = true
)

func main() {

	flag.Parse()

	args := os.Args

	if err := InitConfig(); err != nil {
		redBegin()
		fmt.Printf("找不到配置文件: /etc/gm/conf/gm.toml，详情：" + err.Error() + "\n")
		fmt.Printf("您可能需要运行build_develop.sh、build_staging.sh、build_master.sh其中一个脚本以便初始化gm配置。\n")
		colorEnd()
		goto out_error
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	debugEnv()

	if args == nil || len(args) < 2 {
		usage()
		goto out_error
	}

	switch args[1] {
	case "version":
		fmt.Printf("gm version 1.0.0\n")
		goto out_ok
	case "help":
		usage()
		goto out_ok
	case "install":
		if len(args) != 3 {
			redBegin()
			fmt.Printf("参数错误，命令格式为：gm install id\n")
			colorEnd()
			goto out_error
		}

		id := args[2]
		if err := install(id); err != nil {
			redBegin()
			fmt.Printf("安装失败: %s，正在进行删除操作...\n", err)
			colorEnd()
			if err := remove(id); err != nil {
				redBegin()
				fmt.Printf("删除失败，请手动删除\n")
				colorEnd()
			} else {
				greenBegin()
				fmt.Printf("删除成功！\n")
				colorEnd()
			}
			goto out_error
		} else {
			greenBegin()
			fmt.Printf("安装成功！\n")
			colorEnd()
			goto out_ok
		}
	case "update":
		if len(args) != 3 {
			redBegin()
			fmt.Printf("参数错误，命令格式为：gm update id\n")
			colorEnd()
			goto out_error
		}

		id := args[2]
		if err := update(id); err != nil {
			redBegin()
			fmt.Printf("更新失败: %s\n", err)
			colorEnd()
			goto out_error
		} else {
			greenBegin()
			fmt.Printf("更新成功！\n")
			colorEnd()
			goto out_ok
		}
	case "remove":
		if len(args) != 3 {
			redBegin()
			fmt.Printf("参数错误，命令格式为：gm remove id\n")
			colorEnd()
			goto out_error
		}

		id := args[2]
		if err := remove(id); err != nil {
			redBegin()
			fmt.Printf("删除失败: %s\n", err)
			colorEnd()
			goto out_error
		} else {
			greenBegin()
			fmt.Printf("删除成功！\n")
			colorEnd()
			goto out_ok
		}
	case "list":
		if len(args) != 2 {
			redBegin()
			fmt.Printf("参数错误，命令格式为：gm list\n")
			colorEnd()
			goto out_error
		}

		list()
		goto out_ok
	default:
		redBegin()
		fmt.Printf("不支持该命令，请运行：gm help 查看具体用法\n")
		colorEnd()
		goto out_error
	}

out_ok:
	time.Sleep(100 * time.Millisecond)
	os.Exit(0)

out_error:
	time.Sleep(100 * time.Millisecond)
	os.Exit(1)
}

func usage() {
	fmt.Printf("用法：gm [version] [help] <command> [<args>]\n")
	fmt.Printf("最常用的 gm 命令有：\n")
	fmt.Printf(" 	list 		列出所有项目的状态\n")
	fmt.Printf(" 	install  	安装指定 ID 的项目(比如：安装ID为1的项目 => gm install 1)\n")
	fmt.Printf(" 	update 		更新指定 ID 的项目(比如：更新ID为1的项目 => gm update 1)\n")
	fmt.Printf(" 	remove 		删除指定 ID 的项目(比如：删除ID为1的项目 => gm remove 1)\n")
}

func quit() {
	// TODO
	//log.Debug("Get quit signal")
	//log.Close()
}

func reload() {
	// TODO
	//log.Debug("Get reload signal")
}
