package db

import (
	"strings"
	"time"
	"github.com/chuwt/zing/object"
)

type Balance struct {
	Base      `xorm:"extends"`
	Gateway   string `xorm:"unique(g)"`
	Currency  string `xorm:"unique(g)"`
	Available string // 可用
	Frozen    string // 冻结
}

func (b *Balance) VtBalance() object.VtBalance {
	return object.VtBalance{
		GatewayName: b.Gateway,
		Currency:    b.Currency,
	}
}

func InsertOrUpdateBalance(balance *Balance) error {
	var query = "insert into balance " +
		"(create_time, update_time, gateway, currency, available, frozen) " +
		"values(?,?,?,?,?,?) on duplicate key update "

	update := make([]string, 0)
	if balance.Available != "0" {
		update = append(update, "available="+balance.Available)
	}
	if balance.Frozen != "0" {
		update = append(update, "frozen="+balance.Frozen)
	}
	query += strings.Join(update, ", ")
	_, err := GetEngine().Exec(query,
		time.Now().Unix(),
		time.Now().Unix(),
		balance.Gateway,
		balance.Currency,
		balance.Available,
		balance.Frozen)
	return err
}
