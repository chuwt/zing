package huobi

import (
	"errors"
	"github.com/chuwt/zing/db"
	"github.com/chuwt/zing/object"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
)

type HuoBi struct {
	*Global
	userId    object.UserId
	accountId int64
}

func (g *HuoBi) timeoutConnection(f func()) {
	var step = 2 * time.Second
	for {
		f()
		<-time.After(step)
		step = step * 2
	}
}

func (g *HuoBi) Start(receiver chan object.Event) {
	g.timeoutConnection(func() {
		var err error
	retry:
		if err = g.ws.Init(); err != nil {
			Log.Error("websocket init 失败", zap.Error(err))
			return
		}
		if err = g.ws.Start(receiver); err != nil {
			Log.Error("websocket start 失败", zap.Error(err))
			return
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
				// todo close
				goto retry
			}
		}
	})
}

// 初始化用户gateway
func (g *HuoBi) Init() error {
	var err error
	// 获取用户信息，这里也会检查apiKey
	if err = g.AccountAccounts(); err != nil {
		return err
	}
	//if err := g.AccountAccountsBalance(); err != nil {
	//	return err
	//}
	// 连接ws
	// todo retry 的间隔
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
		//if err := gateway.Factor.AddBalance(g.userId, dbBalance); err != nil {
		//	return err
		//}
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
