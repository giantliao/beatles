package config

import (
	"encoding/json"
	"errors"
	"github.com/kprc/libeth/account"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"sync"
)

const (
	BTL_HomeDir      = ".beatles"
	BTL_CFG_FileName = "beatles.json"
)

type BtlConf struct {
	CmdListenPort string `json:"cmdlistenport"`

	ApiPath           string `json:"api_path"`
	PurchasePath      string `json:"purchase_path"`
	ListMinerPath     string `json:"list_miner_path"`
	RegisterMinerPath string `json:"register_miner_path"`

	HttpServerPort int    `json:"http_server_port"`
	StreamPort     int    `json:"stream_port"`
	StreamIP       net.IP `json:"stream_ip"`

	MasterAccessUrl   string `json:"master_access_url"`
	LicenseServerAddr account.BeatleAddress
	Location          string `json:"location"`

	WalletSavePath string `json:"wallet_save_path"`
	DbPath         string `json:"db_path"`
	ExpireDb       string `json:"expire_db_path"`
}

var (
	btlcfgInst     *BtlConf
	btlcfgInstLock sync.Mutex
)

func (bc *BtlConf) InitCfg() *BtlConf {
	bc.CmdListenPort = "127.0.0.1:50501"
	bc.WalletSavePath = "wallet.json"

	bc.ApiPath = "api"
	bc.PurchasePath = "purchase"
	bc.ListMinerPath = "list"
	bc.RegisterMinerPath = "reg"

	bc.ExpireDb = "expire.db"
	bc.DbPath = "db"

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

func (bc *BtlConf) GetWalletSavePath() string {
	return path.Join(GetBtlHomeDir(), bc.WalletSavePath)
}

func (bc *BtlConf) GetDbPath() string {
	dbpath := path.Join(GetBtlHomeDir(), bc.DbPath)

	if !tools.FileExists(dbpath) {
		os.MkdirAll(dbpath, 0755)
	}

	return dbpath
}

func (bc *BtlConf) GetExpireDbFile() string {
	return path.Join(bc.GetDbPath(), bc.ExpireDb)
}

func IsInitialized() bool {
	if tools.FileExists(GetBtlCFGFile()) {
		return true
	}

	return false
}

func (bc *BtlConf) SetHttpPort(port int) {
	bc.HttpServerPort = port
	bc.Save()
}

func (bc *BtlConf) SetStreamPort(port int) {
	bc.StreamPort = port
	bc.Save()
}

func (bc *BtlConf) SetStreamIP(ipstr string) error {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return errors.New("ip address format not correct")
	}

	bc.StreamIP = ip
	bc.Save()

	return nil
}

func (bc *BtlConf) GetpurchaseWebPath() string {
	return "/" + bc.ApiPath + "/" + bc.PurchasePath
}

func (bc *BtlConf) GetListMinersWebPath() string {
	return "/" + bc.ApiPath + "/" + bc.ListMinerPath
}

func (bc *BtlConf) GetRegisterMinerWebPath() string {
	return "/" + bc.ApiPath + "/" + bc.RegisterMinerPath
}

func (bc *BtlConf) GetMasterAccessUrl() string {
	if strings.Contains(bc.MasterAccessUrl, "http") {
		return bc.MasterAccessUrl
	}

	return "http://" + bc.MasterAccessUrl

}
