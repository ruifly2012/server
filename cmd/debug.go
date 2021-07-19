/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"os"
	"server/gserver"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// debugCmd represents the debug command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "控制台",
	Long:  `gen sever ping GameServer `,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 2 {
			debug(args[0], args[1])
		} else {
			debug(strconv.Itoa(int(gserver.ServerCfg.ServerID)), "127.0.0.1")
		}

	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}

func debug(serverid, ip string) {
	// for i, c := range rootCmd.Commands() {
	// 	fmt.Printf("%v  %v \n", i, c.Name())
	// }
	// rootCmd.Println(rootCmd.UsageString())

	if !connGenServer(serverid, ip) {
		return
	}

	for {
		print("-> ")
		cmd := getInput()

		args := strings.Split(cmd, " ")
		if len(args) > 0 {
			cmd = args[0]
		}

		switch cmd {
		case "getroleinfo":
			call(cmd)
		case "quit":
			return
		case "EOF":
			return
		}
	}
}

func getInput() string {
	//使用os.Stdin开启输入流
	//函数原型 func NewReader(rd io.Reader) *Reader
	//NewReader创建一个具有默认大小缓冲、从r读取的*Reader 结构见官方文档
	in := bufio.NewReader(os.Stdin)
	//in.ReadLine函数具有三个返回值 []byte bool error
	//分别为读取到的信息 是否数据太长导致缓冲区溢出 是否读取失败
	str, _, err := in.ReadLine()
	if err != nil {
		return err.Error()
	}
	return string(str)
}
