package expiredb

import (
	"github.com/giantliao/beatles/config"
	"github.com/kprc/libeth/account"
	"github.com/kprc/nbsnetwork/db"
	"strconv"
	"sync"
)

type ExpireDb struct {
	db.NbsDbInter
	dbLock sync.Mutex
	cursor *db.DBCusor
}


var (
	expireDb *ExpireDb
	expireDbLock sync.Mutex
)

func newExpireDb() *ExpireDb {
	cfg:=config.GetCBtl()
	db:=db.NewFileDb(cfg.GetExpireDbFile()).Load()

	return &ExpireDb{NbsDbInter:db}
}

func GetExpireDb() *ExpireDb  {
	if expireDb == nil{
		expireDbLock.Lock()
		defer expireDbLock.Unlock()
		if expireDb == nil{
			expireDb = newExpireDb()
		}
	}
	return expireDb
}

func (e *ExpireDb) Update(acct account.BeatleAddress,expire int64) {
	e.dbLock.Lock()
	defer e.dbLock.Unlock()

	e.NbsDbInter.Update(acct.String(),strconv.FormatInt(expire,10))

}

func (e *ExpireDb)Find(acct account.BeatleAddress) (int64,error) {
	e.dbLock.Lock()
	defer e.dbLock.Unlock()

	if v, err:=e.NbsDbInter.Find(acct.String());err!=nil{
		return 0,err
	}else {
		var vi int64
		vi, err = strconv.ParseInt(v,10,64)
		if err!=nil{
			return 0,err
		}
		return vi,nil
	}

}



