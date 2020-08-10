package config

import (
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"os"
	"path"
	"sync"
)

const (
	BTL_HomeDir      = ".beatles"
	BTL_CFG_FileName = "beatles.json"
)

type BtlConf struct {
	CmdListenPort  string `json:"cmdlistenport"`
	HttpServerPort int    `json:"http_server_port"`
	WalletSavePath string `json:"wallet_save_path"`

	ApiPath       string `json:"api_path"`
	PurchasePath  string `json:"purchase_path"`
	ListMinerPath string `json:"list_miner_path"`

	StreamPort int    `json:"stream_port"`
	StreamIP   string `json:"stream_ip"`

	MasterAccessUrl string `json:"master_access_url"`
}

var (
	btlcfgInst     *BtlConf
	btlcfgInstLock sync.Mutex
)

func (bc *BtlConf) InitCfg() *BtlConf {
	bc.HttpServerPort = 50511
	bc.CmdListenPort = "127.0.0.1:50501"
	bc.WalletSavePath = "wallet.json"

	bc.ApiPath = "api"
	bc.PurchasePath = "purchase"
	bc.ListMinerPath = "list"

	bc.StreamPort = 50520
	return bc
}

func (bc *BtlConf) Load() *BtlConf {
	if !tools.FileExists(GetBtlCFGFile()) {
		return nil
	}

	jbytes, err := tools.OpenAndReadAll(GetBtlCFGFile())
	if err != nil {
		log.Println("load file failed", err)
		return nil
	}

	err = json.Unmarshal(jbytes, bc)
	if err != nil {
		log.Println("load configuration unmarshal failed", err)
		return nil
	}

	return bc

}

func newBtlmCfg() *BtlConf {

	bc := &BtlConf{}

	bc.InitCfg()

	return bc
}

func GetCBtl() *BtlConf {
	if btlcfgInst == nil {
		btlcfgInstLock.Lock()
		defer btlcfgInstLock.Unlock()
		if btlcfgInst == nil {
			btlcfgInst = newBtlmCfg()
		}
	}

	return btlcfgInst
}

func PreLoad() *BtlConf {
	bc := &BtlConf{}

	return bc.Load()
}

func LoadFromCfgFile(file string) *BtlConf {
	bc := &BtlConf{}

	bc.InitCfg()

	bcontent, err := tools.OpenAndReadAll(file)
	if err != nil {
		log.Fatal("Load Config file failed")
		return nil
	}

	err = json.Unmarshal(bcontent, bc)
	if err != nil {
		log.Fatal("Load Config From json failed")
		return nil
	}

	btlcfgInstLock.Lock()
	defer btlcfgInstLock.Unlock()
	btlcfgInst = bc

	return bc

}

func LoadFromCmd(initfromcmd func(cmdbc *BtlConf) *BtlConf) *BtlConf {
	btlcfgInstLock.Lock()
	defer btlcfgInstLock.Unlock()

	lbc := newBtlmCfg().Load()

	if lbc != nil {
		btlcfgInst = lbc
	} else {
		lbc = newBtlmCfg()
	}

	btlcfgInst = initfromcmd(lbc)

	return btlcfgInst
}

func GetBtlHomeDir() string {
	curHome, err := tools.Home()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(curHome, BTL_HomeDir)
}

func GetBtlCFGFile() string {
	return path.Join(GetBtlHomeDir(), BTL_CFG_FileName)
}

func (bc *BtlConf) Save() {
	jbytes, err := json.MarshalIndent(*bc, " ", "\t")

	if err != nil {
		log.Println("Save BASD Configuration json marshal failed", err)
	}

	if !tools.FileExists(GetBtlHomeDir()) {
		os.MkdirAll(GetBtlHomeDir(), 0755)
	}

	err = tools.Save2File(jbytes, GetBtlCFGFile())
	if err != nil {
		log.Println("Save BASD Configuration to file failed", err)
	}

}

func (bc *BtlConf)GetWalletSavePath() string  {
	return path.Join(GetBtlHomeDir(),bc.WalletSavePath)
}

func (bc *BtlConf) GetPurchasePath() string {
	return "http://" + bc.ApiPath + "/" + bc.PurchasePath
}

func (bc *BtlConf) GetLittMinerPath() string {
	return "http://" + bc.ApiPath + "/" + bc.ListMinerPath
}

func IsInitialized() bool {
	if tools.FileExists(GetBtlCFGFile()) {
		return true
	}

	return false
}
