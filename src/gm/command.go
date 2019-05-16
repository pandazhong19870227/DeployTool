package main

import (
	"bufio"
	//"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

const (
	CommandInstall int = iota
	CommandUpdate
	CommandRemove
)

func install(projectId string) (err error) {
	return runCommand(CommandInstall, projectId)
}

func update(projectId string) (err error) {
	return runCommand(CommandUpdate, projectId)
}

func remove(projectId string) (err error) {
	return runCommand(CommandRemove, projectId)
}

func runCommand(cmdType int, projectId string) (err error) {
	project := getProject(projectId)
	if project == nil {
		return errors.New("找不到该项目，请使用命令：gm list 查看可用项目")
	}

	if !project.Enabled {
		return errors.New("该项目未激活，如需安装，请编辑/etc/gm/conf/gm.toml，设置该项目enabled字段为true")
	}

	scriptPath := ""
	switch cmdType {
	case CommandInstall:
		scriptPath = project.InstallScript
		if ok, _ := project.exists(); ok {
			redBegin()
			fmt.Printf("该项目已安装\n")
			colorEnd()
			return
		}
		break
	case CommandUpdate:
		scriptPath = project.UpdateScript
		if ok, _ := project.exists(); !ok {
			return errors.New("该项目未安装")
		}
		break
	case CommandRemove:
		scriptPath = project.RemoveScript
		if ok, _ := project.exists(); !ok {
			return errors.New("该项目未安装")
		}
		break
	default:
		return errors.New(fmt.Sprintf("无效的命令: %s", cmdType))
	}

	if ok, _ := pathExists(scriptPath); !ok {
		return errors.New(fmt.Sprintf("没有找到对应的shell脚本: %s", scriptPath))
	}

	if !strings.Contains(project.Git, ".git") {
		return errors.New(fmt.Sprintf("无效的git地址: %s，请检查/etc/gm/conf/gm.toml中关于项目%s的配置", project.Git, project.Name))
	}

	if Conf.Env == "" {
		return errors.New("env不能为空，请检查/etc/gm/conf/gm.toml中关于env的配置")
	}

	script := fmt.Sprintf("%s %s %s", scriptPath, project.Git, Conf.Env)

	yellowBegin()
	fmt.Printf("正在执行shell脚本: %s\n", script)
	colorEnd()

	//var out bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", script)
	//cmd.Stdout = &out

	//err = cmd.Run()

	//yellowBegin()
	//fmt.Printf("%s\n", out.String())
	//colorEnd()

	//if err != nil {
	//	err = errors.New(fmt.Sprintf("脚本执行过程中发生错误: %s\n", err))
	//}

	stdout, err2 := cmd.StdoutPipe()

	if err2 != nil {
		err = errors.New(fmt.Sprintf("脚本执行过程中发生错误: %s\n", err2))
		return
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	yellowBegin()
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Printf(line)
	}
	colorEnd()

	if err = cmd.Wait(); err != nil {
		err = errors.New(fmt.Sprintf("脚本执行过程中发生错误: %s\n", err))
	}

	return
}

func getProject(id string) *Project {
	for _, p := range Conf.Projects {
		if p.Id == id {
			return p
		}
	}

	return nil
}

func list() {
	fmt.Printf("已经配置下列项目(项目名称 [项目ID][状态][描述])：\n")
	for _, p := range Conf.Projects {
		/*
			格式：\033[显示方式;前景色;背景色m

			说明：
			前景色            背景色           颜色
			---------------------------------------
			30                40              黑色
			31                41              红色
			32                42              绿色
			33                43              黃色
			34                44              蓝色
			35                45              紫红色
			36                46              青蓝色
			37                47              白色
			显示方式           意义
			-------------------------
			0                终端默认设置
			1                高亮显示
			4                使用下划线
			5                闪烁
			7                反白显示
			8                不可见

			例子：
			\033[1;31;40m
			\033[0m
			如果要加背景色，则
			status = "\033[1;32;40m已安装\033[0m"
		*/
		status := ""
		if !p.Enabled {
			status = "\033[1;33m未激活\033[0m"
		} else {
			ok, _ := p.exists()
			if ok {
				status = "\033[1;32m已安装\033[0m"
			} else {
				status = "\033[1;31m未安装\033[0m"
			}
		}
		fmt.Printf("%30s%5s%5s%5s\n", p.Name, "["+p.Id+"]", "["+status+"]", "\033[1;32m--"+p.Description+"\033[0m")
	}
}

func debugEnv() {
	fmt.Printf("当前是")
	yellowBegin()
	fmt.Printf(" %s ", Conf.Env)
	colorEnd()
	fmt.Printf("环境。\n")
}
