package strategy

import (
	"time"

	"github.com/wiger123/okex_v5_golang/config"
	. "github.com/wiger123/okex_v5_golang/database"
	. "github.com/wiger123/okex_v5_golang/wsdata/client"
)

// PingPong: 每隔一段时间要发起一次交易, 否则 30 秒不操作, 私有频道自动断开连接
func PingPong(c *OkxClient, dr *DataRepo) {
	// 循环
	for {
		// 挂单
		var order1 = c.PostSingleOrder(config.PingPongInstID, config.PingPongTdMode, config.PingPongClOrdId1, config.PingPongSide, config.PingPongPosSide, config.PingPongOrdType, config.PingPongPerSize, config.PingPongPerPrice)
		// 订单
		var order2 = c.PostSingleOrder(config.PingPongInstID, config.PingPongTdMode, config.PingPongClOrdId2, config.PingPongSide, config.PingPongPosSide, config.PingPongOrdType, config.PingPongPerSize, config.PingPongPerPrice)
		// 订单聚合
		var orders []PostOrder
		// 添加订单
		orders = append(orders, order1)
		// 添加订单
		orders = append(orders, order2)
		// 批量下单
		c.PostOrders("posttest", "batch-orders", orders, dr)
		// 等待
		time.Sleep(config.PingPongDelay * time.Second)
		// 订单
		var corder1 = c.CancelSingleOrder(config.PingPongInstID, config.PingPongOrdId1, config.PingPongClOrdId1)
		// 订单
		var corder2 = c.CancelSingleOrder(config.PingPongInstID, config.PingPongOrdId2, config.PingPongClOrdId2)
		// 订单聚合
		var corders []CancelOrder
		// 添加订单
		corders = append(corders, corder1)
		// 添加订单
		corders = append(corders, corder2)
		// 批量撤单
		c.CancelOrders("posttest", "batch-cancel-orders", corders, dr)
		// 等待
		time.Sleep(config.PingPongDelay * time.Second)
	}
}
