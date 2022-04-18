package database

import (
	"log"
	"sync"

	"github.com/wiger123/okex_v5_golang/config"
	. "github.com/wiger123/okex_v5_golang/config"
	. "github.com/wiger123/okex_v5_golang/utils"
	. "github.com/wiger123/okex_v5_golang/wsdata/protocol"
)

// DataRepo 负责数据管理, 接收 websocket 推送的数据, 并对外提供数据获取
type DataRepo struct {
	// 并发保护
	Mu sync.Mutex
	// 交易数据
	TradeData []Trade
	// 盘口数据
	Book5Data []Book5
	// 盘口价格数据
	Book5AvgData []float64
	// 盘口挂买价
	BidsPrice float64
	// 盘口挂卖价
	AsksPrice float64
	// 账户数据
	AccountData []Account
	// Token 数目
	TokenAmt float64
	// USDT 数目
	UsdtAmt float64
	// 持仓数据
	PositionsData []Positions
	// 多仓数据
	PositionsLongData PositionsLong
	// 空仓数据
	PositionsShortData PositionsShort
	// 订单数据
	OrdersData map[string]*Orders
}

// 创建 DataRepo
func NewDataRepo() *DataRepo {
	// 返回结构体
	return &DataRepo{
		// 交易数据
		TradeData: make([]Trade, 0),
		// 盘口数据
		Book5Data: make([]Book5, 0),
		// 盘口价格数据
		Book5AvgData: make([]float64, 0),
		// 盘口挂买价
		BidsPrice: 0,
		// 盘口挂卖价
		AsksPrice: 0,
		// 账户数据
		AccountData: make([]Account, 0),
		// Token 数目
		TokenAmt: 0,
		// USDT 数目
		UsdtAmt: 0,
		// 持仓数据
		PositionsData: make([]Positions, 0),
		// 订单数据
		OrdersData: make(map[string]*Orders),
	}
}

// 处理信息
func (dr *DataRepo) HandleMessage(m PushMessage) {
	// 获取频道名和产品 ID
	channel, _ := m.ChannelAndInstID()
	// 定义错误
	var err error

	// 频道选择
	switch channel {
	// 交易数据
	case "trades":
		// 交易数据信息
		tm := m.(*TradeMessage)
		// 显示数据
		// log.Println("[成功提示] 交易数据: ", tm)
		// 处理数据
		err = dr.handleTrade(tm)
	// 盘口数据
	case "books5":
		// 盘口数据信息
		bm := m.(*Book5Message)
		// 显示数据
		// log.Printf("[成功提示] 盘口数据: %v", bm)
		// 处理数据
		err = dr.handleBook5(bm)
	// 账户数据
	case "account":
		// 账户数据信息
		am := m.(*AccountMessage)
		// 显示数据
		// log.Printf("[成功提示] 账户数据: %v", bm)
		// 处理数据
		err = dr.handleAccount(am)
	// 持仓数据
	case "positions":
		// 持仓数据信息
		pm := m.(*PositionsMessage)
		// 显示数据
		// log.Printf("[成功提示] 持仓数据: %v", pm)
		// 处理数据
		err = dr.handlePositions(pm)
	// 订单数据
	case "orders":
		// 订单数据信息
		om := m.(*OrdersMessage)
		// 显示数据
		// log.Printf("[成功提示] 订单数据: %v", om)
		// 处理数据
		err = dr.handleOrders(om)
	// 未知频道
	default:
		// 普通提示
		log.Printf("[普通提示] 未知频道: %v", channel)
	}

	// 处理信息错误
	if err != nil {
		// 普通提示
		log.Printf("[普通提示] 处理信息出错: %v", err)
	}
}

// 处理交易数据
func (dr *DataRepo) handleTrade(m *TradeMessage) error {
	// 数据库上锁
	dr.Mu.Lock()
	// 函数结束前解锁
	defer dr.Mu.Unlock()
	// 追加数据
	dr.TradeData = append(dr.TradeData, m.Data...)
	// 保留数据
	dr.TradeData = append(dr.TradeData[:0], dr.TradeData[Max(len(dr.TradeData)-Ntrade, 0):]...)
	// 显示数据
	// log.Println("[成功提示] 数据库交易数据: ", dr.TradeData)
	// 显示数据数目
	// log.Println("[成功提示] 数据库交易数据数目: ", len(dr.TradeData))
	// 未出错返回
	return nil
}

// 处理盘口数据
func (dr *DataRepo) handleBook5(m *Book5Message) error {
	// 数据库上锁
	dr.Mu.Lock()
	// 函数结束前解锁
	defer dr.Mu.Unlock()
	// 追加数据
	dr.Book5Data = append(dr.Book5Data, m.Data...)
	// 保留数据
	dr.Book5Data = append(dr.Book5Data[:0], dr.Book5Data[Max(len(dr.Book5Data)-NBook5s, 0):]...)
	// 判断数据是否为空
	if len(m.Data) > 0 {
		// 更新买价
		dr.BidsPrice = 0.618*String2Float64(m.Data[0].Bids[0][0]) + 0.382*String2Float64(m.Data[0].Asks[0][0]) + config.Delta
		// 更新卖价
		dr.AsksPrice = 0.382*String2Float64(m.Data[0].Bids[0][0]) + 0.618*String2Float64(m.Data[0].Asks[0][0]) - config.Delta
		// 盘口加权价格
		avgPrice := (String2Float64(m.Data[0].Asks[0][0])+String2Float64(m.Data[0].Bids[0][0]))*0.35 +
			(String2Float64(m.Data[0].Asks[1][0])+String2Float64(m.Data[0].Bids[1][0]))*0.1 +
			(String2Float64(m.Data[0].Asks[2][0])+String2Float64(m.Data[0].Bids[2][0]))*0.03 +
			(String2Float64(m.Data[0].Asks[3][0])+String2Float64(m.Data[0].Bids[3][0]))*0.015 +
			(String2Float64(m.Data[0].Asks[4][0])+String2Float64(m.Data[0].Bids[4][0]))*0.005
		// 追加数据
		dr.Book5AvgData = append(dr.Book5AvgData, avgPrice)
		// 保留数据
		dr.Book5AvgData = append(dr.Book5AvgData[:0], dr.Book5AvgData[Max(len(dr.Book5AvgData)-NBook5sAvg, 0):]...)
	}
	// 显示数据
	// log.Println("[成功提示] 数据库盘口数据: ", dr.Book5Data)
	// 显示数据数目
	// log.Println("[成功提示] 数据库盘口数据数目: ", len(dr.Book5Data))
	// 未出错返回
	return nil
}

// 处理账户数据
func (dr *DataRepo) handleAccount(m *AccountMessage) error {
	// 数据库上锁
	dr.Mu.Lock()
	// 函数结束前解锁
	defer dr.Mu.Unlock()
	// 追加数据
	dr.AccountData = m.Data
	// 判断数据是否为空
	if len(m.Data) > 0 {
		// 循环
		for i := 0; i < len(m.Data[0].Details); i++ {
			// Token
			if m.Data[0].Details[i].Ccy == config.TokenInstID {
				// 设置余额
				dr.TokenAmt = String2Float64(m.Data[0].Details[i].CashBal)
			}
			// USDT
			if m.Data[0].Details[i].Ccy == config.UsdtInstID {
				// 设置余额
				dr.UsdtAmt = String2Float64(m.Data[0].Details[i].CashBal)
			}
		}

	}
	// 显示数据
	// log.Println("[成功提示] 数据库账户数据: ", dr.AccountData)
	// 未出错返回
	return nil
}

// 处理订单数据: 本地订单簿维护
func (dr *DataRepo) handleOrders(m *OrdersMessage) error {
	// 数据库上锁
	dr.Mu.Lock()
	// 函数结束前解锁
	defer dr.Mu.Unlock()
	// 添加订单
	for i := 0; i < len(m.Data); i++ {
		// 判断 PingPong 订单类型
		if m.Data[i].InstId == config.PingPongInstID {
			// 跳过操作
			continue
		}
		// 初始化
		newOrder := m.Data[i]
		// 挂单
		switch m.Data[i].State {
		// live: 等待成交: 添加到数据库
		case "live":
			// 更新订单簿
			dr.OrdersData[m.Data[i].ClOrdId] = &newOrder
		// partially_filled: 部分成交: 添加到数据库
		case "partially_filled":
			// 更新订单簿
			dr.OrdersData[m.Data[i].ClOrdId] = &newOrder
		// canceled: 撤单成功: 从数据库删除
		case "canceled":
			// 删除元素, 不存在不会报错
			delete(dr.OrdersData, m.Data[i].ClOrdId)
		// filled: 完全成交: 从数据库删除
		case "filled":
			// 删除元素, 不存在不会报错
			delete(dr.OrdersData, m.Data[i].ClOrdId)
		// 其他情况
		default:
			// 暂时不处理
			continue
		}
	}
	// 显示数据
	// log.Println("[成功提示] 数据库订单数据: ", dr.OrdersData)
	// 未出错返回
	return nil
}

// 处理持仓数据
func (dr *DataRepo) handlePositions(m *PositionsMessage) error {
	// 数据库上锁
	dr.Mu.Lock()
	// 函数结束前解锁
	defer dr.Mu.Unlock()
	// 追加数据
	dr.PositionsData = m.Data
	// 数据分类
	for i := 0; i < len(dr.PositionsData); i++ {
		// 多仓归类
		if dr.PositionsData[i].PosSide == "long" {
			/// 持仓方向
			dr.PositionsLongData.PosSide = dr.PositionsData[i].PosSide
			// 持仓数量
			dr.PositionsLongData.Pos = dr.PositionsData[i].Pos
			// 可平仓数量
			dr.PositionsLongData.AvailPos = dr.PositionsData[i].AvailPos
			// 开仓平均价
			dr.PositionsLongData.AvgPx = dr.PositionsData[i].AvgPx
			// 未实现收益
			dr.PositionsLongData.Upl = dr.PositionsData[i].Upl
			// 未实现收益率
			dr.PositionsLongData.UplRatio = dr.PositionsData[i].UplRatio
			// 杠杆倍数
			dr.PositionsLongData.Lever = dr.PositionsData[i].Lever
			// 预估强平价
			dr.PositionsLongData.LiqPx = dr.PositionsData[i].LiqPx
			// 标记价格
			dr.PositionsLongData.MarkPx = dr.PositionsData[i].MarkPx
			// 跳过
			continue
		}
		// 空仓归类
		if dr.PositionsData[i].PosSide == "short" {
			/// 持仓方向
			dr.PositionsShortData.PosSide = dr.PositionsData[i].PosSide
			// 持仓数量
			dr.PositionsShortData.Pos = dr.PositionsData[i].Pos
			// 可平仓数量
			dr.PositionsShortData.AvailPos = dr.PositionsData[i].AvailPos
			// 开仓平均价
			dr.PositionsShortData.AvgPx = dr.PositionsData[i].AvgPx
			// 未实现收益
			dr.PositionsShortData.Upl = dr.PositionsData[i].Upl
			// 未实现收益率
			dr.PositionsShortData.UplRatio = dr.PositionsData[i].UplRatio
			// 杠杆倍数
			dr.PositionsShortData.Lever = dr.PositionsData[i].Lever
			// 预估强平价
			dr.PositionsShortData.LiqPx = dr.PositionsData[i].LiqPx
			// 标记价格
			dr.PositionsShortData.MarkPx = dr.PositionsData[i].MarkPx
			// 跳过
			continue
		}
	}
	// 显示数据
	// log.Println("[成功提示] 数据库持仓数据: ", dr.PositionsData)
	// 未出错返回
	return nil
}
