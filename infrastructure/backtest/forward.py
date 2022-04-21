# 导入依赖
import pandas as pd
import os
import sys
from operationlist import *
from orderlist import *

class Forward:
    """
    回测系统: 构建时间序列, 将未成交的订单保留在订单列表中, 等待后续时间的匹配
    """
    def __init__(self, tickData, start, end, instId, balanceA, balanceB, limitFee, marketFee, delay, level, logInt, pathTrade, pathProfit):
        """
        回测系统初始化
        @param tickData: tick 数据路径
        @param start: 回测起始时间
        @param end: 回测截止时间
        @param instId: 交易品种: A-B
        @param balanceA: A 币初始账户余额
        @param balanceB: B 币初始账户金额
        @param limitFee: 限价手续费
        @param marketFee: 市价手续费
        @param delay: 延迟
        @param level: 盘口深度
        @param logInt: 打印间隔
        @param pathTrade: 交易记录输出路径
        @param pathProfit: 收益记录输出路径
        """
        # 数据库
        self.database = pd.read_csv(tickData)
        # 回测起始时间
        self.start = start
        # 回测截止时间
        self.end = end
        # 交易品种: A-B
        self.instId = instId
        # A 币初始账户余额
        self.balanceA = balanceA
        # B 币初始账户金额
        self.balanceB = balanceB
        # 限价手续费
        self.limitFee = limitFee
        # 市价手续费
        self.marketFee = marketFee
        # 模拟延迟
        self.delay = delay
        # 盘口深度
        self.level = level
        # 打印间隔
        self.logInt = logInt
        # 交易记录输出路径
        self.pathTrade = pathTrade
        # 收益记录输出路径
        self.pathProfit = pathProfit
        # 初始价值
        self.initValue = 0
        # 冻结资金
        self.frozenA = 0
        # 冻结资金
        self.frozenB = 0
        # 交易记录
        self.tradeLog = []
        # 收益记录
        self.profitLog = []
    
    def run(self, operationList):
        """
        启动回测
        """
        # 起始时间
        startInd = self.bisearch(self.start)
        # 结束时间
        endInd = self.bisearch(self.end)
        # 执行回测
        self.backtest(startInd, endInd, operationList)

    def execute(self, time, side, type, size, price):
        """
        执行回调
        @param time: 订单时间: 时间戳
        @param side: 交易方向: buy: 买入 B=>A, sell: 卖出 A=>B
        @param type: 订单类型: LIMIT: 限价单, MARKET: 市价单, CANCEL: 撤销订单
        @param size: 交易数量: A 数目
        @param price: 交易价格: A 价格
        """
        # 限价单
        if type == 'LIMIT':
            # 提示
            print(f"[普通提示] 限价: {'卖出' if side == 'sell' else '买入'} {size} 价格: {price} 时间: {time}")
        # 市价单
        elif type == 'MARKET':
            # 提示
            print(f"[普通提示] 市价: {'卖出' if side == 'sell' else '买入'} {size} 价格: {price} 时间: {time}")
        # 撤销订单
        elif type == 'CANCEL':
            # 提示
            print(f"[普通提示] 撤销: 时间: {time}")

    def logTrade(self):
        """
        交易数据记录
        """

    def logProfit(self):
        """
        收益数据记录
        """

    def calcValue(self, index):
        """
        计算账户价值
        @param index: 指定时间
        """
        # Token A 价值
        valueA = self.database.loc[index, 'bidPx1'] * (self.balanceA + self.frozenA)
        # Token B 价值
        valueB = self.balanceB + self.frozenB
        # 返回价值
        return valueA + valueB

    def backtest(self, startInd, endInd, operationListOrg):
        """
        执行回测
        @param startInd: 回测起始时间索引
        @param endInd: 回测截止时间索引
        @param operationListOrg: 操作列表
        """
        # 账户初始价值
        self.initValue = self.calcValue(0)
        # 交易记录更新
        self.tradeLog = operationListOrg
        # 操作列表
        operationList = operationListOrg
        # 进度打印索引行
        backOver = startInd
        # 订单列表
        orderList = OrderList()
        # 操作列表从第几条开始执行
        operationInd = 0

        # 时间序列
        for index in range(startInd, endInd + 1):
            # 当前时间
            nowTime = self.database.loc[index, 'timeMs']
            # 显示
            # print('[普通提示] 当前时间: {0}'.format(nowTime))

            # 进度打印
            if index >= backOver:
                # 当前账户价值
                nowValue = self.calcValue(index)
                # 收益情况
                profit = nowValue - self.initValue
                # 收益记录更新
                self.profitLog.append(profit)
                # 进度打印
                print('[普通提示] 已回测: {0}%'.format(round((index - startInd) / (endInd + 1 - startInd) * 100, 2)))
                # 索引行更新
                backOver = index + self.logInt

            # 策略列表: 在此时刻有待执行操作, 执行策略
            for opindex in range(operationInd, len(operationList)):
                # 操作时间已超出当前回测时间
                if operationList[opindex].time + self.delay > self.database.loc[index, 'timeMs']:
                    # 下次从这里开始执行
                    operationInd = opindex
                    # 跳出内层循环
                    break
                
                # 操作时间在当前回测时间内
                else:
                    # 下次从下一个开始执行
                    operationInd = opindex + 1
                    # 添加订单
                    if operationList[opindex].operation == 'post':
                        # 买
                        if operationList[opindex].side == 'buy':
                            # 购买后余额
                            newBalanceB = self.balanceB - operationList[opindex].size * operationList[opindex].price
                            # 余额足够
                            if newBalanceB >= 0:
                                # 更新余额
                                self.balanceB = newBalanceB
                                # 更新冻结余额
                                self.frozenB = operationList[opindex].size * operationList[opindex].price
                                # 订单列表: 更新订单列表
                                orderList.post(operationList[opindex])
                                # 提示
                                print('[普通提示] 挂单成功: B Token 账户余额: {0}, B Token 交易消耗: {1}'.format(self.balanceB, operationList[opindex].size * operationList[opindex].price))
                            # 余额不足
                            else:
                                # 提示
                                print('[错误提示] 余额不足: B Token 账户余额: {0}, B Token 交易需要: {1}'.format(self.balanceB, operationList[opindex].size * operationList[opindex].price))

                        # 卖
                        elif operationList[opindex].side == 'sell':
                            # 出售后余额
                            newBalanceA = self.balanceA - operationList[opindex].size
                            # 余额足够
                            if newBalanceA >= 0:
                                # 更新余额
                                self.balanceA = newBalanceA
                                # 更新冻结余额
                                self.frozenA = operationList[opindex].size
                                # 订单列表: 更新订单列表
                                orderList.post(operationList[opindex])
                                # 提示
                                print('[普通提示] 挂单成功: A Token 账户余额: {0}, A Token 交易消耗: {1}'.format(self.balanceA, operationList[opindex].size))
                            # 余额不足
                            else:
                                # 提示
                                print('[错误提示] 余额不足: A Token 账户余额: {0}, A Token 交易需要: {1}'.format(self.balanceA, operationList[opindex].size))
                        # 其他
                        else:
                            # 提示
                            print('[错误提示] 交易类型错误')
                    # 撤销订单
                    elif operationList[opindex].operation == 'cancel':
                        # 判断订单是否已完成, 或者是否已撤销
                        if operationList[opindex].orderId in orderList.orderlist:
                            # 买
                            if orderList.orderlist[operationList[opindex].orderId].side == 'buy':
                                # 归还金额
                                self.balanceB += orderList.orderlist[operationList[opindex].orderId].frozenB
                                # 冻结资金撤销
                                self.frozenB -= orderList.orderlist[operationList[opindex].orderId].frozenB
                                # 订单冻结撤销
                                orderList.orderlist[operationList[opindex].orderId].frozenB = 0
                            # 卖
                            elif orderList.orderlist[operationList[opindex].orderId].side == 'sell':
                                # 归还金额
                                self.balanceA += orderList.orderlist[operationList[opindex].orderId].frozenA
                                # 冻结资金撤销
                                self.frozenA -= orderList.orderlist[operationList[opindex].orderId].frozenA
                                # 订单冻结撤销
                                orderList.orderlist[operationList[opindex].orderId].frozenA = 0

                        # 撤销
                        self.execute(self.database.loc[index, 'timeMs'], '', 'CANCEL', 0, 0)
                        # 订单列表: 更新订单列表
                        orderList.cancel(operationList[opindex].orderId)
                        # 删除本条撤销订单的操作: operation 更改为 over, 以免影响遍历进程
                        operationList[opindex].operation = 'over'
                    # 其他
                    else:
                        # 提示
                        # print('[错误提示] 订单类型错误')
                        pass

            # 订单簿: 在此时刻有待成交订单, 尝试撮合
            for key in list(orderList.orderlist.keys()):
                # 订单
                if orderList.orderlist[key].state == 'full' or orderList.orderlist[key].state == 'cancel':
                    # 跳出
                    continue
                # 买单
                if orderList.orderlist[key].side == 'buy':
                    # 限价单
                    if orderList.orderlist[key].type == 'LIMIT':
                        # 逐档匹配
                        for i in range(1, self.level + 1):
                            # 该档价格
                            px = self.database.loc[index, 'askPx' + str(i)]
                            # 该档数量
                            sz = self.database.loc[index, 'askSz' + str(i)]
                            # 若挂单价格大于该档卖价
                            if orderList.orderlist[key].price >= px:
                                # 待交易数目
                                targetSz = orderList.orderlist[key].size - orderList.orderlist[key].fullSz
                                # 若挂单数目大于该档数目
                                if targetSz > sz:
                                    # 成交该档
                                    self.execute(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, sz, px)
                                    # 订单簿更新
                                    self.database.loc[index, 'askSz' + str(i)] = 0
                                    # 更新余额
                                    self.balanceA += sz
                                    # 已完成数目更新
                                    orderList.orderlist[key].fullSz += sz
                                    # 订单状态更新
                                    orderList.orderlist[key].state = 'part'
                                # 若挂单数目小于等于该档数目
                                else:
                                    # 成交订单
                                    self.execute(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, targetSz, px)
                                    # 订单簿更新
                                    self.database.loc[index, 'askSz' + str(i)] -= targetSz
                                    # 更新余额
                                    self.balanceA += targetSz
                                    # 已完成数目更新
                                    orderList.orderlist[key].fullSz += targetSz
                                    # 订单状态更新
                                    orderList.orderlist[key].state = 'full'
                                    # 订单列表: 更新订单列表
                                    orderList.cancel(orderList.orderlist[key].orderId)
                                    # 跳出循环
                                    break
                            # 若挂单价格小于该档卖价
                            else:
                                # 订单不成交
                                break
                    # 市价单
                    elif orderList.orderlist[key].type == 'MARKET':
                        # 逐档匹配
                        for i in range(1, self.level + 1):
                            # 该档价格
                            px = self.database.loc[index, 'askPx' + str(i)]
                            # 该档数量
                            sz = self.database.loc[index, 'askSz' + str(i)]
                            # 待交易数目
                            targetSz = orderList.orderlist[key].size - orderList.orderlist[key].fullSz
                            # 若挂单数目大于该档数目
                            if targetSz > sz:
                                # 成交该档
                                self.execute(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, sz, px)
                                # 订单簿更新
                                self.database.loc[index, 'askSz' + str(i)] = 0
                                # 更新余额
                                self.balanceA += sz
                                # 已完成数目更新
                                orderList.orderlist[key].fullSz += sz
                                # 订单状态更新
                                orderList.orderlist[key].state = 'part'
                            # 若挂单数目小于等于该档数目
                            else:
                                # 成交订单
                                self.execute(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, targetSz, px)
                                # 订单簿更新
                                self.database.loc[index, 'askSz' + str(i)] -= targetSz
                                # 更新余额
                                self.balanceA += targetSz
                                # 已完成数目更新
                                orderList.orderlist[key].fullSz += targetSz
                                # 订单状态更新
                                orderList.orderlist[key].state = 'full'
                                # 订单列表: 更新订单列表
                                orderList.cancel(orderList.orderlist[key].orderId)
                                # 跳出循环
                                break
                # 卖单
                else:
                    # 限价单
                    if orderList.orderlist[key].type == 'LIMIT':
                        # 逐档匹配
                        for i in range(1, self.level + 1):
                            # 该档价格
                            px = self.database.loc[index, 'bidPx' + str(i)]
                            # 该档数量
                            sz = self.database.loc[index, 'bidSz' + str(i)]
                            # 若挂单价格小于该档买价
                            if orderList.orderlist[key].price <= px:
                                # 待交易数目
                                targetSz = orderList.orderlist[key].size - orderList.orderlist[key].fullSz
                                # 若挂单数目大于该档数目
                                if targetSz > sz:
                                    # 成交该档
                                    self.execute(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, sz, px)
                                    # 订单簿更新
                                    self.database.loc[index, 'bidSz' + str(i)] = 0
                                    # 更新余额
                                    self.balanceB += sz
                                    # 已完成数目更新
                                    orderList.orderlist[key].fullSz += sz
                                    # 订单状态更新
                                    orderList.orderlist[key].state = 'part'
                                # 若挂单数目小于等于该档数目
                                else:
                                    # 成交订单
                                    self.execute(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, targetSz, px)
                                    # 订单簿更新
                                    self.database.loc[index, 'bidSz' + str(i)] -= targetSz
                                    # 更新余额
                                    self.balanceB += sz
                                    # 已完成数目更新
                                    orderList.orderlist[key].fullSz += targetSz
                                    # 订单状态更新
                                    orderList.orderlist[key].state = 'full'
                                    # 订单列表: 更新订单列表
                                    orderList.cancel(orderList.orderlist[key].orderId)
                                    # 跳出循环
                                    break
                            # 若挂单价格大于该档买价
                            else:
                                # 订单不成交
                                break
                    # 市价单
                    elif orderList.orderlist[key].type == 'MARKET':
                        # 逐档匹配
                        for i in range(1, self.level + 1):
                            # 该档价格
                            px = self.database.loc[index, 'bidPx' + str(i)]
                            # 该档数量
                            sz = self.database.loc[index, 'bidSz' + str(i)]
                            # 若挂单数目大于该档数目
                            if orderList.orderlist[key].size > sz:
                                # 成交该档
                                self.execute(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, sz, px)
                                # 订单簿更新
                                self.database.loc[index, 'bidSz' + str(i)] = 0
                                # 更新余额
                                self.balanceB += sz
                                # 已完成数目更新
                                orderList.orderlist[key].fullSz += sz
                                # 订单状态更新
                                orderList.orderlist[key].state = 'part'
                            # 若挂单数目小于等于该档数目
                            else:
                                # 成交订单
                                self.execute(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, targetSz, px)
                                # 订单簿更新
                                self.database.loc[index, 'bidSz' + str(i)] -= targetSz
                                # 更新余额
                                self.balanceB += sz
                                # 已完成数目更新
                                orderList.orderlist[key].fullSz += targetSz
                                # 订单状态更新
                                orderList.orderlist[key].state = 'full'
                                # 订单列表: 更新订单列表
                                orderList.cancel(orderList.orderlist[key].orderId)
                                # 跳出循环
                                break

    def bisearch(self, time):
        """
        二分查找: 获取表格中指定时间的索引
        @param time: 输入时间
        @return high: 索引
        """
        # 起始索引
        low = 0
        # 结束索引
        high = len(self.database) - 1
        # 迭代
        while low <= high:
            # 中值
            mid = low + (high - low) // 2
            # 中值时间大于查找时间
            if self.database.loc[mid, 'timeMs'] > time:
                # 移动上界
                high = mid - 1
            # 中值时间小于查找时间
            else:
                # 移动下界
                low = mid + 1
        # 返回索引
        return high

def _testForward():
    """
    回测测试
    """
    # 当前路径
    modpath = os.path.dirname(os.path.abspath(sys.argv[0]))
    # 数据路径
    tickData = os.path.join(modpath, '../../dataset/tick/20220209DOT-USDT-SWAP.OK.1.csv')
    # 回测起始时间
    start = 1644364800021
    # 回测截止时间
    end = 1644366092038
    # 交易品种: A-B
    instId = 'DOT-USDT'
    # A 币初始账户余额
    balanceA = 5.0
    # B 币初始账户金额
    balanceB = 1000000000000000.0
    # 限价手续费
    limitFee = - 0.025 / 100
    # 市价手续费
    marketFee = 0.03 / 100
    # 模拟延迟
    delay = 100
    # 盘口深度
    level = 5
    # 打印间隔
    logInt = 1000
    # 交易记录输出路径
    pathTrade = os.path.join(modpath, '')
    # 收益记录输出路径
    pathProfit = os.path.join(modpath, '')
    # 初始化
    forward = Forward(tickData, start, end, instId, balanceA, balanceB, limitFee, marketFee, delay, level, logInt, pathTrade, pathProfit)
    
    # 模拟策略操作
    operationList = OperationList()
    # 订单 1
    order1 = Order(1644364800021, 'DOT-USDT', 'buy', 'LIMIT', 2000.0, 21.653, 'test001', 'post')
    # 订单列表更新
    operationList.add(order1)
    # 订单 2
    order2 = Order(1644364990205, 'DOT-USDT', 'buy', 'LIMIT', 10.0, 21.853, 'test001', 'cancel')
    # 订单列表更新
    operationList.add(order2)
    # 订单 3
    order3 = Order(1644364990205, 'DOT-USDT', 'sell', 'LIMIT', 200.0, 21.25, 'test001', 'post')
    # 订单列表更新
    operationList.add(order3)
    
    # 执行回测
    forward.run(operationList.operationList)

    # 显示账户操作历史
    print('[普通提示] 交易历史: {0}'.format(forward.tradeLog))
    # 显示账户收益历史
    print('[普通提示] 收益历史: {0}'.format(forward.profitLog))
    
# 主函数
if __name__ == "__main__":
    # 操作列表测试
    _testForward()