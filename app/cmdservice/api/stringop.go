package api

import (
	"context"
	"github.com/giantliao/beatles/app/cmdcommon"
	"github.com/giantliao/beatles/app/cmdpb"
	"github.com/giantliao/beatles/streamserver"
	"github.com/giantliao/beatles/wallet"
	"github.com/giantliao/beatles/webserver"

	"time"
)

type CmdStringOPSrv struct {
}

func (cso *CmdStringOPSrv) StringOpDo(cxt context.Context, so *cmdpb.StringOP) (*cmdpb.DefaultResp, error) {
	msg := ""
	switch so.Op {
	case cmdcommon.CMD_ACCOUNT_CREATE:
		//msg = createAccount(so.Param[0])
	case cmdcommon.CMD_ACCOUNT_LOAD:
		//msg = loadAccount(so.Param[0])
	case cmdcommon.CMD_RUN:
		//if len(so.)
		msg = run(so.Param[0])
	default:
		return encapResp("Command Not Found"), nil
	}

	return encapResp(msg), nil
}


func run(passwd string) string {


	err := wallet.LoadWallet(passwd)
	if err != nil {
		return "load wallet failed"
	}

	//start web server

	go webserver.StartWebDaemon()

	//start stream server
	go streamserver.StartStreamServer()

	return "start successfully"

}


func int64time2string(t int64) string {
	tm := time.Unix(t/1000, 0)
	return tm.Format("2006-01-02 15:04:05")
}
