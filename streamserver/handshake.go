package streamserver

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/giantliao/beatles-protocol/licenses"
	"github.com/giantliao/beatles-protocol/stream"
	"github.com/giantliao/beatles/config"
	"github.com/giantliao/beatles/expiredb"
	"github.com/giantliao/beatles/wallet"
	"github.com/kprc/libeth/account"
	"github.com/kprc/nbsnetwork/tools"
	"net"
)

func handshake(conn net.Conn) (net.Conn, error) {
	s := &stream.StreamConn{Conn: conn}
	b := stream.NewStreamBuf()

	n, err := s.Read(b)
	if err != nil {
		return nil, err
	}

	var aesk [32]byte
	acct := account.BeatleAddress(b[:n])
	aesk, err = wallet.GetKey(acct)

	cs := stream.NewCipherConn(s, aesk)

	now := tools.GetNowMsTime()
	var expire int64
	expire, err = expiredb.GetExpireDb().Find(acct)
	if err != nil || expire <= now {
		cs.Write([]byte{'1'})
		n, err = cs.Read(b)
		if err != nil {
			return nil, err
		}
		l := &licenses.License{}
		err = json.Unmarshal(b[:n], l)
		if err != nil {
			return nil, err
		}
		//verify license
		if !isValidLicense(acct, l) {
			return nil, errors.New("not a valid signature")
		}

		expiredb.GetExpireDb().Update(acct, l.Content.ExpireTime)
	}
	cs.Write([]byte{'0'})

	return cs, nil
}

func isValidLicense(cid account.BeatleAddress, l *licenses.License) bool {
	if l.Content.Receiver != cid {
		return false
	}

	cfg := config.GetCBtl()

	if l.Content.Provider != cfg.LicenseServerAddr {
		return false
	}

	now := tools.GetNowMsTime()
	if l.Content.ExpireTime < now {
		return false
	}

	forsig, _ := json.Marshal(l.Content)

	return ed25519.Verify(cfg.LicenseServerAddr.DerivePubKey(), forsig, base58.Decode(l.Signature))
}
