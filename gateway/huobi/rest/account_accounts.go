package rest

import (
	"github.com/chuwt/zing/client/http"
)

func (r *HuoBi) AccountAccounts(callback http.Callback) (*AccountAccountsRes, error) {
	short := "/v1/account/accounts"
	path := r.GetUrl(short)
	if callback != nil {
		return nil, nil
	} else {
		signParams := r.api.NewSignParams()
		signature := r.api.Sign("GET", r.Host, short, signParams)
		signParams["Signature"] = signature

		accountAccounts := new(AccountAccountsRes)
		err := http.SyncGetRequest(path, accountAccounts, signParams)
		if err != nil {
			return nil, err
		}
		return accountAccounts, nil
	}
}

type AccountAccountsRes struct {
	BaseRes
	Data []struct {
		Id      int64  `json:"id"`
		Type    string `json:"type"`
		SubType string `json:"subtype"`
		State   string `json:"state"`
	} `json:"data"`
}
