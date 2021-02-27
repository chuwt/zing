# github.com/chuwt/zing

## 简介
基于vnpy的策略平台的go实现，采用go实现底层逻辑，支持多机器和多用户，上层策略兼容vnpy的策略（需要一些小修改）

## 策略说明
- python编写的策略只做指标判断（下单判断）（由于numpy和ta-lib的强大），加载历史数据和订单等放在golang部分
- 后续支持golang原生策略

## 图
![img](https://github.com/chuwt/github.com/chuwt/zing/blob/master/doc/img.png?raw=true)

## todo
- 兼容vnpy策略实现
- 订单/成交处理
- rpc服务
