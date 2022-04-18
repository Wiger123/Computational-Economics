package protocol

// 印钞机属性 print(money)
type PrintMoneyAttr struct {
	// 记录 poll 函数调用时未触发交易的次数, 当触发下单并且下单逻辑执行完时, NumTick 重置为 0
	NumTick int64
	// 交易市场已经成交的订单交易记录时间戳
	LastTradeTime int64
	// 记录当前提单下单后的订单 ID
	TradeOrderId string
	// 通过加权平均计算之后的市场每次考察时成交量参考,
	// 每次循环获取一次市场行情数据, 可以理解为考察了行情一次
	Vol float64
	// 卖单提单价格, 可以理解为策略通过计算后将要挂卖单的价格
	AskPrice float64
	// 买单提单价格
	BidPrice float64
	// 一个数组, 记录订单薄中前三档加权平均计算之后的时间序列上的价格,
	// 简单说就是每次储存计算得到的订单薄前三档加权平均价格, 放在一个数组中,
	// 用于后续策略交易信号参考, 所以该变量名是 prices, 复数形式, 表示一组价格
	Prices []float64
	// 仓位比重, 币的价值正好占总资产价值的一半时, 该值为 0.5, 即平衡状态
	P float64
	// 记录最近一次计算收益时的时间戳, 单位毫秒, 用于控制收益计算部分代码触发执行的频率
	PreCalc int64
	// 记录当前收益数值
	PreNet float64
}

// 创建印钞机
func NewPrintMoney() *PrintMoneyAttr {
	// 返回结构体
	return &PrintMoneyAttr{
		// 策略循环时未触发交易的次数
		NumTick: 0,
		// 上次市场成交时间戳
		LastTradeTime: 0,
		// 自定义订单 id
		TradeOrderId: "",
		// 历史加权成交量
		Vol: 0,
		// 卖单价格
		AskPrice: 0,
		// 买单价格
		BidPrice: 0,
		// 历史盘口价格加权
		Prices: make([]float64, 0),
		// 仓位比重
		P: 0.5,
		// 最近一次计算收益的时间戳
		PreCalc: 0,
		// 当前收益数值
		PreNet: 0,
	}
}
