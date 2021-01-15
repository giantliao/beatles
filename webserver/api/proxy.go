package api

import (
	"fmt"
	"github.com/giantliao/beatles/config"
	"github.com/kprc/nbsnetwork/tools/httputil"
	"io/ioutil"
	"net/http"
	"strings"
)

type BeatlesMasterProxy struct {
}

func (bmp *BeatlesMasterProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "not a post request")
		return
	}
	//read http body error
	if contents, err := ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "read http body error")
		return
	} else {
		proxyUrl := ""
		cfg := config.GetCBtl()
		if strings.Contains(r.URL.Path, cfg.ListMinerPath) {
			proxyUrl = cfg.GetMasterAccessUrl() + cfg.GetListMinersWebPath()
		} else if strings.Contains(r.URL.Path, cfg.PurchasePath) {
			proxyUrl = cfg.GetMasterAccessUrl() + cfg.GetpurchaseWebPath()
		} else if strings.Contains(r.URL.Path, cfg.NoncePrice) {
			proxyUrl = cfg.GetMasterAccessUrl() + cfg.GetNocePriceWebPath()
		} else if strings.Contains(r.URL.Path, cfg.FreshLicensePath) {
			proxyUrl = cfg.GetMasterAccessUrl() + cfg.GetFreshLicensePath()
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "bad rquest url")
			return
		}

		var result string
		var code int
		result, code, err = httputil.Post(proxyUrl, string(contents), false)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, err.Error())
			return
		}

		if code != 200 {
			w.WriteHeader(500)
			fmt.Fprintf(w, "proxy error")
			return
		}

		w.WriteHeader(200)
		w.Write([]byte(result))

		return
	}
}
