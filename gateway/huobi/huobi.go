package huobi

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"vngo/db"
	"vngo/gateway"
	"vngo/object"
)

type HuoBi struct {
	Public
	userId    object.UserId
	accountId int64
}

// 初始化用户gateway
func (g *HuoBi) Init() error {
	// 获取用户信息
	if err := g.AccountAccounts(); err != nil {
		return err
	}
	if err := g.AccountAccountsBalance(); err != nil {
		return err
	}
	// 连接ws
	// todo retry 的间隔
retry:
	if err := g.ws.Init(g.userId); err != nil {
		return err
	}

	for {
		select {
		case err := <-g.ws.Err:
			Log.Error(
				"websocket连接错误",
				zap.String("gateway", string(g.Name())),
				zap.String("userId", string(g.userId)),
				zap.Error(err),
			)
			// todo 重连
			goto retry
		}
	}
}

func (g *HuoBi) Connect() error {
	if err := g.ws.Connect(); err != nil {
		return err
	}
	return nil
}

func (g *HuoBi) AccountAccountsBalance() error {
	accountsBalance, err := g.rest.AccountAccountsBalance(g.accountId, nil)
	if err != nil {
		Log.Error("获取用户余额失败",
			zap.Error(err))
		return err
	}
	for _, balance := range accountsBalance.Data.List {
		if balance.Balance.Equal(decimal.Zero) {
			continue
		}
		// 入库
		dbBalance := &db.Balance{
			Gateway:   string(g.Name()),
			Currency:  balance.Currency,
			Available: "0",
			Frozen:    "0",
		}
		if balance.Type == "trade" {
			dbBalance.Available = balance.Balance.String()
		} else if balance.Type == "frozen" {
			dbBalance.Frozen = balance.Balance.String()
		} else {
			continue
		}
		if err := db.InsertOrUpdateBalance(dbBalance); err != nil {
			Log.Error("更新数据库用户余额失败", zap.Error(err), zap.Any("data", dbBalance))
			return err
		}
		// 同时放入内存
		if err := gateway.Factor.AddBalance(g.userId, dbBalance); err != nil {
			return err
		}
	}
	Log.Info("获取用户余额成功")
	return nil
}

func (g *HuoBi) AccountAccounts() error {
	Log.Info("开始获取accounts")
	accounts, err := g.rest.AccountAccounts(nil)
	if err != nil {
		Log.Error("获取accounts失败",
			zap.Error(err))
		return err
	}
	for _, account := range accounts.Data {
		if account.Type == "spot" {
			g.accountId = account.Id
		}
	}
	if g.accountId == 0 {
		Log.Error("获取accounts失败",
			zap.Any("accounts", accounts))
		return errors.New(" 未找到现货的accountId")
	}
	Log.Info("获取accounts成功", zap.Int64("accountId", g.accountId))
	return nil
}
