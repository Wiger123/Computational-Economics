package config

// 参数配置
const (
	// 交易品种
	InstID = "DOGE-USDT"
	// 品种类型
	InstType = "SPOT"
	// 交易模式
	TdMode = "cash"
	// 挂单模式
	OrdType = "limit"
	// 小数点后几位
	FloatPrec = 2
	// 订单编号长度
	ClOrdIdLength = 10
	// 交易量强弱比
	Ratio = 3.0
	// 开单最低交易量
	MinTradeVolume = 200
	// 止盈交易量强弱比
	CoverRatio = 2.0
	// 止盈最低交易量
	CoverMinTradeVolume = 50
	// 档位划分
	NumLevel = 10
	// 最高档位权重
	MaxWeight = 1.0
	// 最低档位权重
	MinWeight = 0.1
	// 档位划分
	NumPost = 10
	// 最大挂单量
	MaxPost = 1.0
	// 最低张数
	MinPost = 1.0
	// 参考权重上限百分比
	MaxRef = 10000.0
	// 卖单挂盘口档位
	AsksLevel = 0
	// 买单挂盘口档位
	BidsLevel = 0
	// 平空仓买单档位
	CoverShortLevel = 0
	// 平多仓卖单档位
	CoverLongLevel = 0
	// 撤单定时 Millisecond
	TimeCancel = 2000
	// 止盈百分比
	StopProfit = 0.05
	// 止损百分比
	StopLoss = -0.05
	// 杠杆倍数
	Leverage = 3
)
