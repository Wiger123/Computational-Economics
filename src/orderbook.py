import pandas as pd
import os.path
import sys

class Order:
    def __init__(self, time, side, size, price):
        """
        Class representing a single order
        @param time: order time 
        @param side: 0=buy, 1=sell
        @param size: order quantity
        @param price: price in ticks
        """
        self.time = time
        self.side = side
        self.size = size
        self.price = price
        self.state = None # ['Completed', 'Canceled']...

class OrderBook(object):
    def __init__(self, name, file_path, cb=None):
        """
        Creates a new order book. 
        
        @param name: name of asset being traded (e.g "BTCUSD")
        @param file_path: order book filepath
        @param cb: trade execution callback. Should have same signature as self.execute
        """
        self.name = name
        self.database = pd.read_csv(file_path)
        self.cb = cb or self.execute

    def execute(self, time, side, type, price, size):
        """
        Execution callback
        @param trader_buy: the trader on the buy side
        @param trader_sell: the trader on the sell side
        @param price: trade price
        @param size: trade size
        """
        if type == 'LIMIT':
            print(f"LIMIT: {'SELL' if side == 1 else 'BUY'} {size} PRICE {price} at TIME {time}")
        elif type == 'MARKET':
            print(f"MARKET: {'SELL' if side == 1 else 'BUY'} {size} PRICE {price} at TIME {time}")
        # print("EXECUTE: %s BUY %s SELL %s %s @ %d" % (trader_buy, trader_sell, size, self.name, price))
    
    def order(self, time, side, type, price, size):
        """
        Match Orders

        @param time: order time
        @param side: 0=buy, 1=sell
        @param type: ['LIMIT', 'MARKET]
        @param price: order price
        @param size: order size
        """
        index = self.bisearch(time)
        # item = self.database.iloc[index]
        if side == 0: # buy order
            if type == 'LIMIT': # limit order
                for i in range(1, 6):
                    px = self.database.loc[index, 'askPx' + str(i)] #askPx
                    sz = self.database.loc[index, 'askSz' + str(i)] #askSz
                    if price >= px:
                        if size > sz:
                            self.cb(time, side, type, px, sz)
                            self.database.loc[index, 'askSz' + str(i)] = 0
                            size -= sz
                        else:
                            self.cb(time, side, type, px, size)
                            self.database.loc[index, 'askSz' + str(i)] -= size
                            break
                    else:
                        # order incomplete
                        break
            elif type == 'MARKET': # market order
                for i in range(1, 6):
                    px = self.database.loc[index, 'askPx' + str(i)] #askPx
                    sz = self.database.loc[index, 'askSz' + str(i)] #askSz
                    if size > sz:
                        self.cb(time, side, type, px, sz)
                        self.database.loc[index, 'askSz' + str(i)] = 0
                        size -= sz
                    else:
                        self.cb(time, side, type, px, size)
                        self.database.loc[index, 'askSz' + str(i)] -= size
                        break
                # if size > 0: order incomplete
            return
        else: # sell order
            if type == 'LIMIT': # limit order
                for i in range(1, 6):
                    px = self.database.loc[index, 'bidPx' + str(i)]
                    sz = self.database.loc[index, 'bidSz' + str(i)]
                    if price <= px:
                        if size > sz:
                            self.cb(time, side, type, px, sz)
                            self.database.loc[index, 'bidSz' + str(i)] = 0
                            size -= sz
                        else:
                            self.cb(time, side, type, px, size)
                            self.database.loc[index, 'bidSz' + str(i)] -= size
                            break
                    else:
                        # order incomplete
                        break                
            elif type == 'MARKET': # market order
                for i in range(1, 6):
                    px = self.database.loc[index, 'bidPx' + str(i)] #bidPx
                    sz = self.database.loc[index, 'bidSz' + str(i)] #bidSz
                    if size > sz:
                        self.cb(time, side, type, px, sz)
                        self.database.loc[index, 'bidSz' + str(i)] = 0
                        size -= sz
                    else:
                        self.cb(time, side, type, px, size)
                        self.database.loc[index, 'bidSz' + str(i)] -= size
                        break
                # if size > 0: order incomplete
            return

    def bisearch(self, time):
        low = 0
        high = len(self.database) - 1
        while low <= high:
            mid = low + (high - low) // 2
            if self.database.loc[mid, 'timeMs'] > time:
                high = mid - 1
            else:
                low = mid + 1
        return high

def _unittest1():
    # Datas are in a subfolder of the samples. Need to find where the script is
    # because it could have been called from anywhere
    modpath = os.path.dirname(os.path.abspath(sys.argv[0]))
    datapath = os.path.join(modpath, '../data/20220209DOT-USDT-SWAP.OK.1.csv')

    book = OrderBook("TEST", datapath)
    
    # time, side, type, price, size
    book.order(1644364800021, 0, 'LIMIT', 21.653, 10.0)
    book.order(1644364800101, 1, 'LIMIT', 21.651, 80.0)
    book.order(1644364800101, 1, 'MARKET', 21.647, 80.0)


if __name__ == "__main__":
    _unittest1()

#TODO: 目前是直接读取了整个OrderBook离线操作的，只实现了交易的逻辑，后续也许可以转换为在线回测
#TODO: 是否需要对于每个Order建立一个对象，从而传递一些信息，比如交易状态，交易时间等
#TODO: 代码比较丑陋，性能还有提升的空间