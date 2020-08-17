// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"errors"
	"fmt"
	"github.com/giantliao/beatles/register"
	"github.com/giantliao/beatles/streamserver"
	"github.com/giantliao/beatles/wallet"
	"github.com/giantliao/beatles/webserver"
	"github.com/howeyc/gopass"
	"github.com/kprc/libeth/account"
	"net"
	"os"

	"github.com/giantliao/beatles/app/cmdcommon"
	"github.com/giantliao/beatles/config"

	"github.com/giantliao/beatles/app/cmdservice"

	"github.com/spf13/cobra"
	"log"
)

////var cfgFile string
//
var (
	cmdconfigfilename string
)

var keypassword string

func inputpassword() (password string, err error) {
	passwd, err := gopass.GetPasswdPrompt("Please Enter Password: ", true, os.Stdin, os.Stdout)
	if err != nil {
		return "", err
	}

	if len(passwd) < 1 {
		return "", errors.New("Please input valid password")
	}

	return string(passwd), nil
}

func inputChoose() (choose string, err error) {
	c, err := gopass.GetPasswdPrompt("Do you reinit config[yes/no]: ", true, os.Stdin, os.Stdout)
	if err != nil {
		return "", err
	}

	return string(c), nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cc",
	Short: "start beatles in current shell",
	Long:  `start beatles in current shell`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := cmdcommon.IsProcessCanStarted()
		if err != nil {
			log.Println(err)
			return
		}

		InitCfg()
		cfg := config.GetCBtl()

		if cfg.LicenseServerAddr == "" || cfg.Location == "" || cfg.MasterAccessUrl == "" {
			log.Println("please initial first")
			return
		}

		cfg.Save()

		if keypassword == "" {
			if keypassword, err = inputpassword(); err != nil {
				log.Println(err)
				return
			}
		}

		err = wallet.LoadWallet(keypassword)
		if err != nil {
			panic("load wallet failed")
		}

		log.Println("register self to beatles master")
		err = register.RegMiner()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("start web server")
		go webserver.StartWebDaemon()

		log.Println("start stream server")
		go streamserver.StartStreamServer()

		cmdservice.GetCmdServerInst().StartCmdService()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func InitCfg() {
	if cmdconfigfilename != "" {
		cfg := config.LoadFromCfgFile(cmdconfigfilename)
		if cfg == nil {
			return
		}
	} else {
		config.LoadFromCmd(cfginit)
	}

}

func cfginit(bc *config.BtlConf) *config.BtlConf {
	cfg := bc
	if masterAccessUrl != "" {
		cfg.MasterAccessUrl = masterAccessUrl
	}
	if licenseServerBetalesAddr != "" {
		cfg.LicenseServerAddr = account.BeatleAddress(licenseServerBetalesAddr)
	}
	if minerLocation != "" {
		cfg.Location = minerLocation
	}
	if minerServerAddr != "" {
		cfg.StreamIP = net.ParseIP(minerServerAddr)
	}

	return cfg

}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringVarP(&cmdconfigfilename, "config-file-name", "c", "", "configuration file name")

}
