# 导入依赖
from order import *

class OrderList:
    """
    订单列表: 将未成交的订单保留在订单列表中, 等待后续时间的匹配
    """
    def __init__(self):
        """
        订单列表初始化
        """
        # 订单列表
        self.orderlist = {}
    
    def post(self, order):
        """
        添加订单
        @param order: 订单
        """
        # key-value 形式
        self.orderlist[order.orderId] = order
        # 普通提示
        # print('[普通提示] 订单添加成功')

    def cancel(self, orderId):
        """
        撤销订单
        @param orderId: 订单编号
        """
        # 撤销失败: 已完全成交, 或未挂单成功
        if orderId not in self.orderlist:
            # 普通提示
            print('[普通提示] 订单撤销失败 {0}: 该订单不存在, 或已完全成交'.format(orderId))
        # 撤销成功: 部分成交, 或未成交
        else:
            # 删除该订单
            self.orderlist.pop(orderId)
            # 普通提示
            print('[普通提示] 订单消失 {0}: 该订单撤销成功, 或该订单交易成功'.format(orderId))


def _testOrderList():
    """
    订单列表测试
    """
    # 初始化
    orderList = OrderList()
    # 订单 1
    order1 = Order(1644364800021, 'DOT-USDT', 'buy', 'LIMIT', 10.0, 21.653, 'test001', 'post')
    # 订单列表更新
    orderList.post(order1)
    # 显示订单列表
    print('[普通提示] 订单列表: {0}'.format(orderList.orderlist))
    # 订单 2
    order2 = Order(1644364800205, 'DOT-USDT', 'buy', 'LIMIT', 10.0, 21.853, 'test002', 'post')
    # 订单列表更新
    orderList.post(order2)
    # 显示订单列表
    print('[普通提示] 订单列表: {0}'.format(orderList.orderlist))
    # 撤销订单
    orderList.cancel(order1.orderId)
    # 显示订单列表
    print('[普通提示] 订单列表: {0}'.format(orderList.orderlist))
    # 撤销订单
    orderList.cancel(order2.orderId)
    # 显示订单列表
    print('[普通提示] 订单列表: {0}'.format(orderList.orderlist))

# 主函数
if __name__ == "__main__":
    # 订单列表测试
    _testOrderList()