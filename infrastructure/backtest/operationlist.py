# 导入依赖
from order import *

class OperationList:
    """
    操作列表: 整个回测过程中所有通过策略的执行输出
    """
    def __init__(self):
        """
        操作列表初始化
        """
        # 订单列表
        self.operationList = []
    
    def add(self, order):
        """
        添加操作
        @param order: 订单
        """
        # 列表形式
        self.operationList.append(order)
        # 普通提示
        print('[普通提示] 操作添加成功')

def _testOperationList():
    """
    订单列表测试
    """
    # 初始化
    operationList = OperationList()
    # 订单 1
    order1 = Order(1644364800021, 'DOT-USDT', 'buy', 'LIMIT', 10.0, 21.653, 'test001', 'post')
    # 订单列表更新
    operationList.add(order1)
    # 显示订单列表
    print('[普通提示] 订单列表: {0}'.format(operationList.operationList))
    # 订单 2
    order2 = Order(1644364800205, 'DOT-USDT', 'buy', 'LIMIT', 10.0, 21.853, 'test002', 'cancel')
    # 订单列表更新
    operationList.add(order2)
    # 显示订单列表
    print('[普通提示] 订单列表: {0}'.format(operationList.operationList))
    
# 主函数
if __name__ == "__main__":
    # 操作列表测试
    _testOperationList()