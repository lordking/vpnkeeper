// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"vpnkeeper/vpn"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "vpnkeeper",
	Short: "Mac平台最方便的VPN管理工具",
	Run: func(cmd *cobra.Command, args []string) {
		stepOne()
		stepTwo()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default .vpnkeeper.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".vpnkeeper")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func readInput() string {
	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		str := string(data)
		if str != "" {
			return str
		}
	}

	return ""
}

//第一步，选择VPN
func stepOne() {

	services, err := vpn.Fetch()

	if err != nil {
		fmt.Print("Find vpn error:", err)
		os.Exit(0)
	}

	if len(services) == 0 {
		fmt.Print("Not found vpn.\n")
		os.Exit(0)
	}

	n := 1
	for _, service := range services[0:] {
		fmt.Printf("%d) %s\n", n, service.Name)
		n++
	}

	selected := viper.GetString("selected")
	if selected != "" {
		seq, _ := strconv.Atoi(selected)
		seq = seq - 1
		vpn.Select(seq)

	} else {
		fmt.Printf("Above is the ordinal sort of VPN in your macOS. Please enter a VPN serial number[%d-%d]: ", 1, len(services))
		selected = readInput()
		seq, _ := strconv.Atoi(selected)
		seq = seq - 1
		err := vpn.Select(seq)
		if err != nil {
			fmt.Print("\nSorry! You input a wrong number! Please again!\n")
			stepOne()
		}

	}
}

func stepTwo() {

	go vpn.RunServ()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case killSignal := <-interrupt:
			if killSignal == os.Interrupt {
				vpn.Stop()
				fmt.Print("Service was interruped by system signal.\n")
			}
			fmt.Print("Service was killed.\n")

			time.Sleep(1 * time.Second)
			return
		}
	}
}
