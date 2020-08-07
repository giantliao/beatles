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
	BTLM_HomeDir      = ".btlmaster"
	BTLM_CFG_FileName = "btlmaster.json"
)

type BtlMasterConf struct {
	EthAccessPoint string `json:"eth_access_point"`
	TrxAccessPoint string `json:"trx_access_point"`

	CmdListenPort  string `json:"cmdlistenport"`
	HttpServerPort int    `json:"http_server_port"`
	WalletSavePath string `json:"wallet_save_path"`

	ApiPath       string `json:"api_path"`
	PurchasePath  string `json:"purchase_path"`
	ListMinerPath string `json:"list_miner_path"`
}

var (
	btlmcfgInst     *BtlMasterConf
	btmlcfgInstLock sync.Mutex
)

func (bc *BtlMasterConf) InitCfg() *BtlMasterConf {
	bc.HttpServerPort = 50101
	bc.CmdListenPort = "127.0.0.1:50500"
	bc.WalletSavePath = "wallet.json"

	bc.ApiPath = "api"
	bc.PurchasePath = "purchase"
	bc.ListMinerPath = "list"

	return bc
}

func (bc *BtlMasterConf) Load() *BtlMasterConf {
	if !tools.FileExists(GetBtlmCFGFile()) {
		return nil
	}

	jbytes, err := tools.OpenAndReadAll(GetBtlmCFGFile())
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

func newBtlmCfg() *BtlMasterConf {

	bc := &BtlMasterConf{}

	bc.InitCfg()

	return bc
}

func GetCBtlm() *BtlMasterConf {
	if btlmcfgInst == nil {
		btmlcfgInstLock.Lock()
		defer btmlcfgInstLock.Unlock()
		if btlmcfgInst == nil {
			btlmcfgInst = newBtlmCfg()
		}
	}

	return btlmcfgInst
}

func PreLoad() *BtlMasterConf {
	bc := &BtlMasterConf{}

	return bc.Load()
}

func LoadFromCfgFile(file string) *BtlMasterConf {
	bc := &BtlMasterConf{}

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

	btmlcfgInstLock.Lock()
	defer btmlcfgInstLock.Unlock()
	btlmcfgInst = bc

	return bc

}

func LoadFromCmd(initfromcmd func(cmdbc *BtlMasterConf) *BtlMasterConf) *BtlMasterConf {
	btmlcfgInstLock.Lock()
	defer btmlcfgInstLock.Unlock()

	lbc := newBtlmCfg().Load()

	if lbc != nil {
		btlmcfgInst = lbc
	} else {
		lbc = newBtlmCfg()
	}

	btlmcfgInst = initfromcmd(lbc)

	return btlmcfgInst
}

func GetBtlmCHomeDir() string {
	curHome, err := tools.Home()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(curHome, BTLM_HomeDir)
}

func GetBtlmCFGFile() string {
	return path.Join(GetBtlmCHomeDir(), BTLM_CFG_FileName)
}

func (bc *BtlMasterConf) Save() {
	jbytes, err := json.MarshalIndent(*bc, " ", "\t")

	if err != nil {
		log.Println("Save BASD Configuration json marshal failed", err)
	}

	if !tools.FileExists(GetBtlmCHomeDir()) {
		os.MkdirAll(GetBtlmCHomeDir(), 0755)
	}

	err = tools.Save2File(jbytes, GetBtlmCFGFile())
	if err != nil {
		log.Println("Save BASD Configuration to file failed", err)
	}

}

func (bc *BtlMasterConf) GetPurchasePath() string {
	return "http://" + bc.ApiPath + "/" + bc.PurchasePath
}

func (bc *BtlMasterConf) GetLittMinerPath() string {
	return "http://" + bc.ApiPath + "/" + bc.ListMinerPath
}

func IsInitialized() bool {
	if tools.FileExists(GetBtlmCFGFile()) {
		return true
	}

	return false
}
