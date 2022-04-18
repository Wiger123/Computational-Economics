package strategy

// import (
// 	"time"

// 	"github.com/wiger123/okex_v5_golang/config"
// 	. "github.com/wiger123/okex_v5_golang/wsdata/client"
// )

// // 执行策略
// func OnStrategy(c *OkxClient) {
// 	// 挂单
// 	var order1 = c.PostSingleOrder(config.InstID, config.TdMode, "20220328130100AE86", "buy", "long", "post_only", "30", "5")
// 	// 订单
// 	var order2 = c.PostSingleOrder(config.InstID, config.TdMode, "20220328130100M78X", "buy", "long", "post_only", "30", "6")
// 	// 订单聚合
// 	var orders []PostOrder
// 	// 添加订单
// 	orders = append(orders, order1)
// 	// 添加订单
// 	orders = append(orders, order2)
// 	// 批量下单
// 	c.PostOrders("posttest", "batch-orders", orders)
// 	// 等待
// 	time.Sleep(5 * time.Second)
// 	// 订单
// 	var corder1 = c.CancelSingleOrder(config.InstID, "", "20220328130100AE86")
// 	// 订单
// 	var corder2 = c.CancelSingleOrder(config.InstID, "", "20220328130100M78X")
// 	// 订单聚合
// 	var corders []CancelOrder
// 	// 添加订单
// 	corders = append(corders, corder1)
// 	// 添加订单
// 	corders = append(corders, corder2)
// 	// 批量撤单
// 	c.CancelOrders("posttest", "batch-cancel-orders", corders)
// }
