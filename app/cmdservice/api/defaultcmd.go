package api

import (
	"context"
	"encoding/json"
	"github.com/giantliao/beatles/app/cmdcommon"
	"github.com/giantliao/beatles/app/cmdpb"
	"github.com/giantliao/beatles/config"
	"github.com/giantliao/beatles/wallet"
	"time"
)

type CmdDefaultServer struct {
	Stop func()
}

func (cds *CmdDefaultServer) DefaultCmdDo(ctx context.Context,
	request *cmdpb.DefaultRequest) (*cmdpb.DefaultResp, error) {

	msg := ""

	switch request.Reqid {
	case cmdcommon.CMD_STOP:
		msg = cds.stop()
	case cmdcommon.CMD_CONFIG_SHOW:
		msg = cds.configShow()
	case cmdcommon.CMD_WALLET_SHOW:
		msg = cds.showWallet()
	}

	if msg == "" {
		msg = "No Results"
	}

	resp := &cmdpb.DefaultResp{}
	resp.Message = msg

	return resp, nil

}

func (cds *CmdDefaultServer) stop() string {

	go func() {
		time.Sleep(time.Second * 2)
		cds.Stop()
	}()

	return "beatles stopped"
}

func encapResp(msg string) *cmdpb.DefaultResp {
	resp := &cmdpb.DefaultResp{}
	resp.Message = msg

	return resp
}

func (cds *CmdDefaultServer) configShow() string {
	cfg := config.GetCBtl()

	bapc, err := json.MarshalIndent(*cfg, "", "\t")
	if err != nil {
		return "Internal error"
	}

	return string(bapc)
}

func (cds *CmdDefaultServer) showWallet() string {
	if _, err := wallet.GetWallet(); err != nil {
		return err.Error()
	} else {
		var s string
		if s, err = wallet.ShowWallet(); err != nil {
			return err.Error()
		}

		return s
	}
}
