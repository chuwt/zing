package rest

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/chuwt/zing/client/http"
)

func (r *HuoBi) AccountAccountsBalance(accountId int64, callback http.Callback) (*AccountAccountsBalanceRes, error) {
	short := fmt.Sprintf("/v1/account/accounts/%d/balance", accountId)
	path := r.GetUrl(short)
	if callback != nil {
		return nil, nil
	} else {
		signParams := r.api.NewSignParams()
		signature := r.api.Sign("GET", r.Host, short, signParams)
		signParams["Signature"] = signature

		accountAccountsBalanceRes := new(AccountAccountsBalanceRes)
		err := http.SyncGetRequest(path, accountAccountsBalanceRes, signParams)
		if err != nil {
			return nil, err
		}
		return accountAccountsBalanceRes, nil
	}
}

type AccountAccountsBalanceRes struct {
	Data struct {
		List []struct {
			Balance  decimal.Decimal `json:"balance"`  //	余额
			Currency string          `json:"currency"` //	币种
			Type     string          `json:"type"`     //	trade: 交易余额，frozen: 冻结余额, loan: 待还借贷本金, interest: 待还借贷利息, lock: 锁仓, bank: 储蓄
		} `json:"list"`
	} `json:"data"`
}
