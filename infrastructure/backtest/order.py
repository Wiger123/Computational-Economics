class Order:
    """
    订单: 单个订单, 及相关操作
    """
    def __init__(self, time, instId, side, type, size, price, orderId, operation):
        """
        订单初始化
        @param time: 订单时间: 时间戳
        @param instId: 交易品种: A-B
        @param side: 交易方向: buy: 买入 B=>A, sell: 卖出 A=>B
        @param type: 订单类型: LIMIT: 限价单, MARKET: 市价单
        @param size: 交易数量: A 数目
        @param price: 交易价格: A 价格
        @param fullSz: 成交数量: 初始为 0 后续更新
        @param orderId: 订单编号: 用于查询, 撤销订单
        @param operation: 订单操作: post: 添加订单, cancel: 撤销订单
        """
        # 订单时间
        self.time = time
        # 交易品种
        self.instId = instId
        # 交易方向
        self.side = side
        # 订单类型
        self.type = type
        # 交易数量
        self.size = size
        # 交易价格
        self.price = price
        # 订单状态: full, part, wait, cancel
        self.state = 'wait'
        # 成交量
        self.fullSz = 0
        # 订单编号
        self.orderId = orderId
        # 操作
        self.operation = operation
        # 冻结 Token A: 在发单时设置为冻结的数目, 撤销订单, 交易完成时归还剩余数目
        self.frozenA = 0
        # 冻结 Token B: 在发单时设置为冻结的数目, 撤销订单, 交易完成时归还剩余数目
        self.frozenB = 0