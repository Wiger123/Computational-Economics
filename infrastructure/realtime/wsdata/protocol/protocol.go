package protocol

// 所有 websocket 推送信息接口
type PushMessage interface {
	// 定义接口
	ChannelAndInstID() (string, string)
}

// 交易数据
type Trade struct {
	// 产品 ID
	InstID string `json:"instId"`
	// 成交 ID
	TradeID string `json:"tradeId"`
	// 成交价格
	Px string `json:"px"`
	// 成交数量
	Sz string `json:"sz"`
	// 成交方向
	Side string `json:"side"`
	// 成交时间
	Ts string `json:"ts"`
}

// 交易数据推送信息
type TradeMessage struct {
	// 频道参数
	Arg struct {
		// 频道名
		Channel string `json:"channel"`
		// 产品 ID
		InstID string `json:"instId"`
	} `json:"arg"`
	// 交易数据
	Data []Trade `json:"data"`
}

// 解析交易数据
func (tm *TradeMessage) ChannelAndInstID() (string, string) {
	// 返回频道名称, 产品 ID
	return tm.Arg.Channel, tm.Arg.InstID
}

// 盘口数据
type Book5 struct {
	// 卖方深度
	Asks [][]string `json:"asks"`
	// 买方深度
	Bids [][]string `json:"bids"`
	// 时间戳
	Ts string `json:"ts"`
	// 校验和
	Checksum int `json:"checksum"`
}

// 盘口数据推送信息
type Book5Message struct {
	// 频道参数
	Arg struct {
		// 频道名
		Channel string `json:"channel"`
		// 产品 ID
		InstID string `json:"instId"`
	} `json:"arg"`
	// 增量 or 全量推送数据
	Action string `json:"action"`
	// 盘口数据
	Data []Book5 `json:"data"`
}

// 解析盘口数据
func (bm *Book5Message) ChannelAndInstID() (string, string) {
	// 返回频道名称, 产品 ID
	return bm.Arg.Channel, bm.Arg.InstID
}

// 账户资产详情
type AccountDetails struct {
	// 币种
	Ccy string `json:"ccy"`
	// 币种总权益
	Eq string `json:"eq"`
	// 币种余额
	CashBal string `json:"cashBal"`
	// 币种余额信息的更新时间
	UTime string `json:"uTime"`
	// 币种逐仓仓位权益
	IsoEq string `json:"isoEq"`
	// 可用保证金
	AvailEq string `json:"availEq"`
	// 美金层面币种折算权益
	DisEq string `json:"disEq"`
	// 可用余额
	AvailBal string `json:"availBal"`
	// 币种占用金额
	FrozenBal string `json:"frozenBal"`
	// 挂单冻结数量
	OrdFrozen string `json:"ordFrozen"`
	// 币种负债额
	Liab string `json:"liab"`
	// 未实现盈亏
	Upl string `json:"upl"`
	// 由于仓位未实现亏损导致的负债
	UplLiab string `json:"uplLiab"`
	// 币种全仓负债额
	CrossLiab string `json:"crossLiab"`
	// 币种逐仓负债额
	IsoLiab string `json:"isoLiab"`
	// 保证金率
	MgnRatio string `json:"mgnRatio"`
	// 计息
	Interest string `json:"interest"`
	// 当前负债币种触发系统自动换币的风险
	Twap string `json:"twap"`
	// 币种最大可借
	MaxLoan string `json:"maxLoan"`
	// 币种权益美金价值
	EqUsd string `json:"eqUsd"`
	// 币种杠杆倍数
	NotionalLever string `json:"notionalLever"`
	// 币种美元指数
	CoinUsdPrice string `json:"coinUsdPrice"`
	// 策略权益
	StgyEq string `json:"stgyEq"`
	// 逐仓未实现盈亏
	IsoUpl string `json:"isoUpl"`
}

// 账户数据
type Account struct {
	// 获取账户信息的最新时间
	UTime string `json:"uTime"`
	// 美金层面权益
	TotalEq string `json:"totalEq"`
	// 美金层面逐仓仓位权益
	IsoEq string `json:"isoEq"`
	// 美金层面有效保证金
	AdjEq string `json:"adjEq"`
	// 美金层面全仓挂单占用保证金
	OrdFroz string `json:"ordFroz"`
	// 美金层面占用保证金
	Imr string `json:"imr"`
	// 美金层面维持保证金
	Mmr string `json:"mmr"`
	// 美金层面保证金率
	MgnRatio string `json:"mgnRatio"`
	// 以美金价值为单位的持仓数量
	NotionalUsd string `json:"notionalUsd"`
	// 各币种资产详细信息
	Details []AccountDetails `json:"details"`
}

// 账户数据推送信息
type AccountMessage struct {
	// 频道参数
	Arg struct {
		// 频道名
		Channel string `json:"channel"`
		// 用户标识
		Uid string `json:"uid"`
		// 币种
		Ccy string `json:"ccy"`
	} `json:"arg"`
	// 账户数据
	Data []Account `json:"data"`
}

// 解析盘口数据
func (bm *AccountMessage) ChannelAndInstID() (string, string) {
	// 返回频道名称, 产品 ID
	return bm.Arg.Channel, ""
}

// 持仓数据
type Positions struct {
	// 产品类型
	InstType string `json:"instType"`
	// 保证金模式
	MgnMode string `json:"mgnMode"`
	// 持仓 ID
	PosId string `json:"posId"`
	// 持仓方向
	PosSide string `json:"posSide"`
	// 持仓数量
	Pos string `json:"pos"`
	// 交易币余额
	BaseBal string `json:"baseBal"`
	// 计价币余额
	QuoteBal string `json:"quoteBal"`
	// 持仓数量币种
	PosCcy string `json:"posCcy"`
	// 可平仓数量
	AvailPos string `json:"availPos"`
	// 开仓平均价
	AvgPx string `json:"avgPx"`
	// 未实现收益
	Upl string `json:"upl"`
	// 未实现收益率
	UplRatio string `json:"uplRatio"`
	// 产品 ID
	InstId string `json:"instId"`
	// 杠杆倍数
	Lever string `json:"lever"`
	// 预估强平价
	LiqPx string `json:"liqPx"`
	// 标记价格
	MarkPx string `json:"markPx"`
	// 初始保证金
	Imr string `json:"imr"`
	// 保证金余额
	Margin string `json:"margin"`
	// 保证金率
	MgnRatio string `json:"mgnRatio"`
	// 维持保证金
	Mmr string `json:"mmr"`
	// 负债额
	Liab string `json:"liab"`
	// 负债币种
	LiabCcy string `json:"liabCcy"`
	// 利息
	Interest string `json:"interest"`
	// 最新成交 ID
	TradeId string `json:"tradeId"`
	// 以美金价值为单位的持仓数量
	NotionalUsd string `json:"notionalUsd"`
	// 期权价值
	OptVal string `json:"optVal"`
	// 信号区
	Adl string `json:"adl"`
	// 占用保证金的币种
	Ccy string `json:"ccy"`
	// 最新成交价
	Last string `json:"last"`
	// 美金价格
	UsdPx string `json:"usdPx"`
	// 美金本位持仓仓位 delta
	DeltaBS string `json:"deltaBS"`
	//  币本位持仓仓位 delta
	DeltaPA string `json:"deltaPA"`
	//  美金本位持仓仓位 gamma
	GammaBS string `json:"gammaBS"`
	//  币本位持仓仓位 gamma
	GammaPA string `json:"gammaPA"`
	//  美金本位持仓仓位 theta
	ThetaBS string `json:"thetaBS"`
	//  币本位持仓仓位 theta
	ThetaPA string `json:"thetaPA"`
	//  美金本位持仓仓位 vega
	VegaBS string `json:"vegaBS"`
	//  币本位持仓仓位 vega
	VegaPA string `json:"vegaPA"`
	//  持仓创建时间
	CTime string `json:"cTime"`
	// 最近一次持仓更新时间
	UTime string `json:"uTime"`
	// 持仓信息的推送时间
	PTime string `json:"pTime"`
}

// 持仓数据推送信息
type PositionsMessage struct {
	// 频道参数
	Arg struct {
		// 频道名
		Channel string `json:"channel"`
		// 用户标识
		Uid string `json:"uid"`
		// 产品类型
		InstType string `json:"instType"`
		// 标的指数
		Uly string `json:"uly"`
		// 产品 ID
		InstId string `json:"instId"`
	} `json:"arg"`
	// 持仓数据
	Data []Positions `json:"data"`
}

// 解析订单数据
func (pm *PositionsMessage) ChannelAndInstID() (string, string) {
	// 返回频道名称, 产品 ID
	return pm.Arg.Channel, pm.Arg.InstId
}

// 多仓空仓数据
type PositionsLong struct {
	// 持仓方向
	PosSide string `json:"posSide"`
	// 持仓数量
	Pos string `json:"pos"`
	// 可平仓数量
	AvailPos string `json:"availPos"`
	// 开仓平均价
	AvgPx string `json:"avgPx"`
	// 未实现收益
	Upl string `json:"upl"`
	// 未实现收益率
	UplRatio string `json:"uplRatio"`
	// 杠杆倍数
	Lever string `json:"lever"`
	// 预估强平价
	LiqPx string `json:"liqPx"`
	// 标记价格
	MarkPx string `json:"markPx"`
}

// 多仓空仓数据
type PositionsShort struct {
	// 持仓方向
	PosSide string `json:"posSide"`
	// 持仓数量
	Pos string `json:"pos"`
	// 可平仓数量
	AvailPos string `json:"availPos"`
	// 开仓平均价
	AvgPx string `json:"avgPx"`
	// 未实现收益
	Upl string `json:"upl"`
	// 未实现收益率
	UplRatio string `json:"uplRatio"`
	// 杠杆倍数
	Lever string `json:"lever"`
	// 预估强平价
	LiqPx string `json:"liqPx"`
	// 标记价格
	MarkPx string `json:"markPx"`
}

// 订单数据
type Orders struct {
	// 产品类型
	InstType string `json:"instType"`
	// 产品ID
	InstId string `json:"instId"`
	// 保证金币种
	Ccy string `json:"ccy"`
	// 订单ID
	OrdId string `json:"ordId"`
	// 由用户设置的订单 ID 来识别您的订单
	ClOrdId string `json:"clOrdId"`
	// 订单标签
	Tag string `json:"tag"`
	// 委托价格
	Px string `json:"px"`
	// 原始委托数量
	Sz string `json:"sz"`
	// 委托单预估美元价值
	NotionalUsd string `json:"notionalUsd"`
	// 订单类型
	OrdType string `json:"ordType"`
	// 订单方向
	Side string `json:"side"`
	// 持仓方向
	PosSide string `json:"posSide"`
	// 交易模式
	TdMode string `json:"tdMode"`
	// 市价单委托数量的类型
	TgtCcy string `json:"tgtCcy"`
	// 最新成交价格
	FillPx string `json:"fillPx"`
	// 最新成交 ID
	TradeId string `json:"tradeId"`
	// 最新成交数量
	FillSz string `json:"fillSz"`
	// 最新成交时间
	FillTime string `json:"fillTime"`
	// 最新一笔成交的手续费
	FillFee string `json:"fillFee"`
	// 最新一笔成交的手续费币种
	FillFeeCcy string `json:"fillFeeCcy"`
	// 最新一笔成交的流动性方向
	ExecType string `json:"execType"`
	// 累计成交数量
	AccFillSz string `json:"accFillSz"`
	// 委托单已成交的美元价值
	FillNotionalUsd string `json:"fillNotionalUsd"`
	// 成交均价
	AvgPx string `json:"avgPx"`
	// 订单状态
	State string `json:"state"`
	// 杠杆倍数
	Lever string `json:"lever"`
	// 止盈触发价
	TpTriggerPx string `json:"tpTriggerPx"`
	// 止盈触发价类型
	TpTriggerPxType string `json:"tpTriggerPxType"`
	// 止盈委托价
	TpOrdPx string `json:"tpOrdPx"`
	// 止损触发价
	SlTriggerPx string `json:"slTriggerPx"`
	// 止损触发价类型
	SlTriggerPxType string `json:"slTriggerPxType"`
	// 止损委托价
	SlOrdPx string `json:"slOrdPx"`
	// 交易手续费币种
	FeeCcy string `json:"feeCcy"`
	// 订单交易手续费
	Fee string `json:"fee"`
	// 返佣金币种
	RebateCcy string `json:"rebateCcy"`
	// 返佣金额
	Rebate string `json:"rebate"`
	// 收益
	Pnl string `json:"pnl"`
	// 订单来源
	Source string `json:"source"`
	// 订单种类分类
	Category string `json:"category"`
	// 订单更新时间
	UTime string `json:"uTime"`
	// 订单创建时间
	CTime string `json:"cTime"`
	// 修改订单时使用的 request ID
	ReqId string `json:"reqId"`
	// 修改订单的结果
	AmendResult string `json:"amendResult"`
	// 错误码
	Code string `json:"code"`
	// 错误消息
	Msg string `json:"msg"`
}

// 订单信息
type OrdersMessage struct {
	// 请求订阅的频道列表
	Arg struct {
		// 频道名
		Channel string `json:"channel"`
		// 用户标识
		Uid string `json:"uid"`
		// 产品类型
		InstType string `json:"instType"`
		// 标的指数
		Uly string `json:"uly"`
		// 产品 ID
		InstId string `json:"instId"`
	} `json:"arg"`
	// 订阅的数据
	Data []Orders `json:"data"`
}

// 解析订单数据
func (om *OrdersMessage) ChannelAndInstID() (string, string) {
	// 返回频道名称, 产品 ID
	return om.Arg.Channel, om.Arg.InstId
}
