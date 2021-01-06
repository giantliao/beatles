package register

import (
	"errors"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/giantliao/beatles-protocol/miners"
	"github.com/giantliao/beatles/config"
	"github.com/giantliao/beatles/port"
	"github.com/giantliao/beatles/wallet"
	"github.com/kprc/nbsnetwork/tools/httputil"
	"log"
	"strconv"
	"time"
)

var regKeepAliveChan chan struct{}

func RegMiner() error {
	m := &miners.Miner{}

	cfg := config.GetCBtl()

	if cfg.StreamPort == 0 {
		cfg.StreamPort = port.TcpPort()
	}

	w, err := wallet.GetWallet()
	if err != nil {
		return err
	}

	m.MinerId = w.BtlAddress()
	m.Location = cfg.Location
	m.Port = cfg.StreamPort
	if cfg.StreamIP != nil {
		m.Ipv4Addr = cfg.StreamIP.String()
	}

	var aesk []byte
	aesk, err = w.AesKey2(cfg.LicenseServerAddr)
	if err != nil {
		return err
	}

	var content []byte
	content, err = m.Marshal(aesk)
	if err != nil {
		return err
	}

	mt := &meta.Meta{}

	mt.Content = content
	mt.Marshal(w.BtlAddress().String(), content)

	regUrl := cfg.GetMasterAccessUrl() + cfg.GetRegisterMinerWebPath()
	var result string
	var code int
	result, code, err = httputil.Post(regUrl, mt.ContentS, false)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New("register failed, code : " + strconv.Itoa(code))
	}

	log.Println("register miner self ", result)

	return nil
}

func StartKeepAlive()  {

	regKeepAliveChan = make(chan struct{},1)

	tic:=time.NewTicker(time.Second*300)
	defer tic.Stop()

	for{
		select {
		case <-tic.C:
			RegMiner()
		case <-regKeepAliveChan:
			return
		}
	}
}

func StopKeepAlive()  {
	close(regKeepAliveChan)
}
