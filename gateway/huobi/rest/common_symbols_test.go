package rest

import "testing"

func TestCommonSymbols(t *testing.T) {
	r := HuoBi{
		Host: "https://api.huobipro.com",
	}
	res, err := r.CommonSymbols(nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
