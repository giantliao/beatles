/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/giantliao/beatles/app/cmdcommon"
	"github.com/giantliao/beatles/config"
	"github.com/kprc/libeth/account"
	"net"

	"github.com/spf13/cobra"
	"log"
)

var masterAccessUrl string
var licenseServerBetalesAddr string
var minerLocation string
var minerServerAddr string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init beatles",
	Long:  `init beatles`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		_, err = cmdcommon.IsProcessCanStarted()
		if err != nil {
			log.Println(err)
			return
		}
		if masterAccessUrl == "" || net.ParseIP(masterAccessUrl) == nil {
			log.Println("please set correct master ip")
			return
		}

		if licenseServerBetalesAddr == "" || !(account.BeatleAddress(licenseServerBetalesAddr).IsValid()) {
			log.Println("please beatles address")
			return
		}

		if minerLocation == "" {
			log.Println("please miner location")
			return
		}
		if minerServerAddr != "" && net.ParseIP(minerServerAddr) == nil {
			log.Println("please set correct miner ip address")
			return
		}

		InitCfg()

		cfg := config.GetCBtl()

		cfg.Save()

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	//initCmd.Flags().StringVarP(&keypassword, "password", "p", "", "password for key encrypt")
	initCmd.Flags().StringVarP(&masterAccessUrl, "master-ip", "m", "", "master ip address")
	initCmd.Flags().StringVarP(&licenseServerBetalesAddr, "master-beatles-addr", "b", "", "beatles address")
	initCmd.Flags().StringVarP(&minerLocation, "miner-geography-addr", "g", "", "geography address")
	initCmd.Flags().StringVarP(&minerServerAddr, "local-stream-server-ip", "s", "", "miner stream server ip address")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
