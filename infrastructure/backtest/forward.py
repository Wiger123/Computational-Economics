# 导入依赖
import pandas as pd
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
    
    def run(self):
        """
        启动回测
        """
        # 起始时间
        startInd = self.bisearch(self.start)
        # 结束时间
        endInd = self.bisearch(self.end)
        # 执行回测
        self.backtest(startInd, endInd)

    def execute(self, time, side, type, size, price):
        """
        执行回调
        @param time: 订单时间: 时间戳
        @param side: 交易方向: buy: 买入 B=>A, sell: 卖出 A=>B
        @param type: 订单类型: LIMIT: 限价单, MARKET: 市价单
        @param size: 交易数量: A 数目
        @param price: 交易价格: A 价格
        """
        # 限价单
        if type == 'LIMIT':
            # 提示
            print(f"限价: {'卖出' if side == 'sell' else '买入'} {size} 价格: {price} 时间: {time}")
        # 市价单
        elif type == 'MARKET':
            # 提示
            print(f"市价: {'卖出' if side == 'sell' else '买入'} {size} 价格: {price} 时间: {time}")

    def backtest(self, startInd, endInd, operationList):
        """
        执行回测
        @param startInd: 回测起始时间索引
        @param endInd: 回测截止时间索引
        @param operationList: 操作列表
        """
        # 进度打印索引行
        backOver = startInd
        # 订单列表
        orderList = OrderList()
        # 操作列表从第几条开始执行
        operationInd = 0

        # 时间序列
        for index in range(startInd, endInd + 1):
            # 进度打印
            if index >= backOver:
                # 进度打印
                print('[普通提示] 已回测: {0}%'.format(round((index - startInd) / (endInd + 1 - startInd) * 100, 2)))
                # 索引行更新
                backOver = index + self.logInt
                
            # 策略列表: 在此时刻有待执行操作, 执行策略
            for opindex in range(operationInd, len(operationList)):
                # 操作时间已超出当前回测时间
                if operationInd[opindex].time + self.delay > self.database.loc[index, 'timeMs']:
                    # 下次从这里开始执行
                    operationInd = opindex
                    # 跳出内层循环
                    break
                
                # 操作时间在当前回测时间内
                else:
                    # 添加订单
                    if operationInd[opindex].operation == 'post':
                        # 买
                        if operationInd[opindex].side == 'buy':
                            # 购买后余额
                            newBalanceB = self.balanceB - operationInd[opindex].size * operationInd[opindex].price
                            # 余额足够
                            if newBalanceB >= 0:
                                # 更新余额
                                self.balanceB = newBalanceB
                                # 订单列表: 更新订单列表
                                orderList.post(operationInd[opindex])
                                # 提示
                                print('[普通提示] 挂单成功: B Token 账户余额: {0}, B Token 交易消耗: {1}'.format(self.balanceB, operationInd[opindex].size * operationInd[opindex].price))
                            # 余额不足
                            else:
                                # 提示
                                print('[错误提示] 余额不足: B Token 账户余额: {0}, B Token 交易需要: {1}'.format(self.balanceB, operationInd[opindex].size * operationInd[opindex].price))

                        # 卖
                        elif operationInd[opindex].side == 'sell':
                            # 出售后余额
                            newBalanceA = self.balanceA - operationInd[opindex].size
                            # 余额足够
                            if newBalanceA >= 0:
                                # 更新余额
                                self.balanceA = newBalanceA
                                # 订单列表: 更新订单列表
                                orderList.post(operationInd[opindex])
                                # 提示
                                print('[普通提示] 挂单成功: A Token 账户余额: {0}, A Token 交易消耗: {1}'.format(self.balanceA, operationInd[opindex].size))
                            # 余额不足
                            else:
                                # 提示
                                print('[错误提示] 余额不足: A Token 账户余额: {0}, A Token 交易需要: {1}'.format(self.balanceA, operationInd[opindex].size))
                        # 其他
                        else:
                            # 提示
                            print('[错误提示] 交易类型错误')
                    # 撤销订单
                    elif operationInd[opindex].operation == 'cancel':
                        # 订单列表: 更新订单列表
                        orderList.cancel(operationInd[opindex])
                    # 其他
                    else:
                        # 提示
                        print('[错误提示] 订单类型错误')

            # 订单簿: 在此时刻有待成交订单, 尝试撮合
            for key in orderList:
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
                                # 若挂单数目大于该档数目
                                if orderList.orderlist[key].size > sz:
                                    # 成交该档
                                    self.cb(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, sz, px)
                                    # 订单簿更新
                                    self.database.loc[index, 'askSz' + str(i)] = 0
                                    # 订单列表更新
                                    orderList.orderlist[key].size -= sz
                                    # 更新余额
                                    self.balanceA += sz
                                # 若挂单数目小于等于该档数目
                                else:
                                    # 成交订单
                                    self.cb(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, orderList.orderlist[key].size, px)
                                    # 订单簿更新
                                    self.database.loc[index, 'askSz' + str(i)] -= orderList.orderlist[key].size
                                    # 更新余额
                                    self.balanceA += orderList.orderlist[key].size
                                    # 跳出循环
                                    break
                            # 若挂单价格小于该档卖价
                            else:
                                # 订单不成交
                                break
                    # 市价单
                    elif type == 'MARKET':
                        # 逐档匹配
                        for i in range(1, self.level + 1):
                            # 该档价格
                            px = self.database.loc[index, 'askPx' + str(i)]
                            # 该档数量
                            sz = self.database.loc[index, 'askSz' + str(i)]
                            # 若挂单数目大于该档数目
                            if orderList.orderlist[key].size > sz:
                                # 成交该档
                                self.cb(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, sz, px)
                                # 订单簿更新
                                self.database.loc[index, 'askSz' + str(i)] = 0
                                # 订单列表更新
                                orderList.orderlist[key].size -= sz
                                # 更新余额
                                self.balanceA += sz
                            # 若挂单数目小于等于该档数目
                            else:
                                # 成交订单
                                self.cb(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, orderList.orderlist[key].size, px)
                                # 订单簿更新
                                self.database.loc[index, 'askSz' + str(i)] -= orderList.orderlist[key].size
                                # 更新余额
                                self.balanceA += orderList.orderlist[key].size
                                # 跳出循环
                                break
                # 卖单
                else:
                    # 限价单
                    if type == 'LIMIT':
                        # 逐档匹配
                        for i in range(1, self.level + 1):
                            # 该档价格
                            px = self.database.loc[index, 'bidPx' + str(i)]
                            # 该档数量
                            sz = self.database.loc[index, 'bidSz' + str(i)]
                            # 若挂单价格小于该档买价
                            if orderList.orderlist[key].price <= px:
                                # 若挂单数目大于该档数目
                                if orderList.orderlist[key].size > sz:
                                    # 成交该档
                                    self.cb(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, sz, px)
                                    # 订单簿更新
                                    self.database.loc[index, 'bidSz' + str(i)] = 0
                                    # 订单列表更新
                                    orderList.orderlist[key].size -= sz
                                    # 更新余额
                                    self.balanceB += sz
                                # 若挂单数目小于等于该档数目
                                else:
                                    # 成交订单
                                    self.cb(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, orderList.orderlist[key].size, px)
                                    # 订单簿更新
                                    self.database.loc[index, 'bidSz' + str(i)] -= orderList.orderlist[key].size
                                    # 更新余额
                                    self.balanceB += sz
                                    # 跳出循环
                                    break
                            # 若挂单价格大于该档买价
                            else:
                                # 订单不成交
                                break
                    # 市价单
                    elif type == 'MARKET':
                        # 逐档匹配
                        for i in range(1, self.level + 1):
                            # 该档价格
                            px = self.database.loc[index, 'bidPx' + str(i)]
                            # 该档数量
                            sz = self.database.loc[index, 'bidSz' + str(i)]
                            # 若挂单数目大于该档数目
                            if orderList.orderlist[key].size > sz:
                                # 成交该档
                                self.cb(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, sz, px)
                                # 订单簿更新
                                self.database.loc[index, 'bidSz' + str(i)] = 0
                                # 订单列表更新
                                orderList.orderlist[key].size -= sz
                                # 更新余额
                                self.balanceB += sz
                            # 若挂单数目小于等于该档数目
                            else:
                                # 成交订单
                                self.cb(self.database.loc[index, 'timeMs'], orderList.orderlist[key].side, orderList.orderlist[key].type, orderList.orderlist[key].size, px)
                                # 订单簿更新
                                self.database.loc[index, 'bidSz' + str(i)] -= orderList.orderlist[key].size
                                # 更新余额
                                self.balanceB += sz
                                # 跳出循环
                                break

    def logTrade(self):
        """
        交易数据记录
        """

    def logProfit(self):
        """
        收益数据记录
        """

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