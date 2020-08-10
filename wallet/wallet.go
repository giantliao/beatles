package wallet


import (
	"errors"
	"github.com/giantliao/beatles/config"
	"github.com/kprc/libeth/wallet"
	"github.com/kprc/nbsnetwork/tools"
)

var (
	beatlesWallet wallet.WalletIntf
)

func GetWallet() (wallet.WalletIntf, error) {
	if beatlesWallet == nil {
		return nil, errors.New("no wallet, please load it")
	}
	return beatlesWallet, nil
}

func newWallet(auth, savepath, remoteeth string) wallet.WalletIntf {
	w := wallet.CreateWallet(savepath, remoteeth)

	if w == nil {
		return nil
	}

	w.Save(auth)

	return w
}

func LoadWallet(auth string) error {
	cfg := config.GetCBtl()

	if !tools.FileExists(cfg.GetWalletSavePath()) {
		beatlesWallet = newWallet(auth, cfg.GetWalletSavePath(), "")
		if beatlesWallet == nil {
			return errors.New("create wallet failed ")
		}
	} else {
		var err error
		beatlesWallet, err = wallet.RecoverWallet(cfg.GetWalletSavePath(), "", auth)
		if err != nil {
			return errors.New("load wallet failed : " + err.Error())
		}
	}
	return nil
}







