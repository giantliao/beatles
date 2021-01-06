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
	"github.com/giantliao/beatles/app/cmdcommon"
	"github.com/giantliao/beatles/app/cmdservice"
	"github.com/giantliao/beatles/config"
	"github.com/kprc/nbsnetwork/tools/processChan"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"log"
	"path"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "beatles start in backend",
	Long:  `beatles start in backend`,
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

		daemondir := config.GetBtlHomeDir()
		cntxt := daemon.Context{
			PidFileName: path.Join(daemondir, "beatles.pid"),
			PidFilePerm: 0644,
			LogFileName: path.Join(daemondir, "beatles.log"),
			LogFilePerm: 0640,
			WorkDir:     daemondir,
			Umask:       027,
			Args:        []string{},
		}
		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			if keypassword == "" {
				if keypassword, err = inputpassword(); err != nil {
					log.Println(err)
					return
				}
			}

			processChan.SendPasswd(daemondir,keypassword)
			log.Println("beatles starting, please check log at:", path.Join(daemondir, "beatles.log"))
			return
		}
		defer cntxt.Release()

		passwd:=processChan.ReceivePasswd(daemondir)

		start(passwd)

		cmdservice.GetCmdServerInst().StartCmdService()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// daemonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	daemonCmd.Flags().StringVarP(&cmdconfigfilename, "config-file-name", "c", "", "configuration file name")

}
