package wallet

import (
	"github.com/kprc/libeth/account"
	"sync"
)

var (
	keymaplock sync.Mutex
	keymap     map[account.BeatleAddress][32]byte
	emptykey   [32]byte
)

func init() {
	keymap = make(map[account.BeatleAddress][32]byte)
}

func GetKey(acct account.BeatleAddress) ([32]byte, error) {
	if v, ok := keymap[acct]; ok {
		return v, nil
	}
	keymaplock.Lock()
	defer keymaplock.Unlock()
	if v, ok := keymap[acct]; ok {
		return v, nil
	}

	w, err := GetWallet()
	if err != nil {
		return emptykey, err
	}

	var (
		key  []byte
		aesk [32]byte
	)

	key, err = w.AesKey2(acct)
	if err != nil {
		return emptykey, err
	}
	copy(aesk[:], key)

	return aesk, nil
}
