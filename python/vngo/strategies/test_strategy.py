# coding:utf-8

from template import (
    CtaTemplate,
    StopOrder,
    TickData,
    BarData,
    TradeData,
    OrderData
)

from time import time


class TestStrategy(CtaTemplate):
    """"""
    author = "用Python的交易员"

    test_trigger = 10

    tick_count = 0
    test_all_done = False

    parameters = ["test_trigger"]
    variables = ["tick_count", "test_all_done"]

    def __init__(self, cta_engine, strategy_name, vt_symbol, setting):
        """"""
        super().__init__(cta_engine, strategy_name, vt_symbol, setting)

        self.test_funcs = [
            self.test_market_order,
            self.test_limit_order,
            self.test_cancel_all,
            self.test_stop_order
        ]
        self.last_tick = None
        self.contract = None

    def on_init(self):
        """
        Callback when strategy is inited.
        """
        return self.success("策略初始化")

    def on_start(self):
        """
        Callback when strategy is started.
        """
        return self.success("策略启动")

    def on_stop(self):
        """
        Callback when strategy is stopped.
        """
        self.success("策略停止")

    def on_tick(self, tick: TickData):
        """
        Callback of new tick data update.
        """
        if self.test_all_done:
            return self.success("测试全部完成")

        self.last_tick = tick

        self.tick_count += 1
        if self.tick_count >= self.test_trigger:
            self.tick_count = 0

            if self.test_funcs:
                test_func = self.test_funcs.pop(0)

                start = time()
                test_func()
                time_cost = (time() - start) * 1000
                return self.success("耗时%s毫秒" % (time_cost))
            else:
                self.test_all_done = True
                return self.success("测试全部完成")

    def on_bar(self, bar: BarData):
        """
        Callback of new bar data update.
        """
        pass

    def on_order(self, order: OrderData):
        """
        Callback of new order data update.
        """
        self.put_event()

    def on_trade(self, trade: TradeData):
        """
        Callback of new trade data update.
        """
        self.put_event()

    def on_stop_order(self, stop_order: StopOrder):
        """
        Callback of stop order update.
        """
        self.put_event()

    def test_market_order(self):
        """"""
        # self.buy(self.last_tick.limit_up, 1)
        return self.success("执行市价单测试")

    def test_limit_order(self):
        """"""
        # self.buy(self.last_tick.limit_down, 1)
        return self.success("执行限价单测试")

    def test_stop_order(self):
        """"""
        # self.buy(self.last_tick.ask_price_1, 1, True)
        return self.success("执行停止单测试")

    def test_cancel_all(self):
        """"""
        # self.cancel_all()
        return self.success("执行全部撤单测试")