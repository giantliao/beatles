package api

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/giantliao/beatles-protocol/meta"
	"github.com/giantliao/beatles/wallet"
	"github.com/kprc/libeth/account"
	w2 "github.com/kprc/libeth/wallet"
	"log"
	"io/ioutil"
	"net/http"
)

type PingTest struct {

}

func (pt *PingTest)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method,r.URL.Path)
	_, _, _, wal, err := DecodeMeta(r)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	buf:=make([]byte,32)
	_,err = rand.Read(buf)
	if err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	resp := &meta.Meta{}
	resp.Marshal(wal.BtlAddress().String(), buf)

	w.WriteHeader(200)
	fmt.Fprint(w, resp.ContentS)
}


func DecodeMeta(r *http.Request) (key []byte, cipherTxt []byte, sender string, w w2.WalletIntf, err error) {
	if r.Method != "POST" {
		return nil, nil, "", nil, errors.New("not a post request")
	}
	var body []byte

	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, nil, "", nil, errors.New("read http body error")
	}

	req := &meta.Meta{ContentS: string(body)}

	sender, cipherTxt, err = req.UnMarshal()
	if err != nil || !(account.BeatleAddress(sender).IsValid()) {
		return nil, nil, "", nil, errors.New("not a correct request")
	}

	w, err = wallet.GetWallet()
	if err != nil {
		return nil, nil, "", nil, errors.New("server have no wallet")
	}

	key, err = w.AesKey2(account.BeatleAddress(sender))
	if err != nil {
		return nil, nil, "", nil, err
	}

	return
}
