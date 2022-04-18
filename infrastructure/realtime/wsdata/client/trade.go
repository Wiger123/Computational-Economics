package client

import (
	"log"

	"github.com/wiger123/okex_v5_golang/config"
	. "github.com/wiger123/okex_v5_golang/database"
	. "github.com/wiger123/okex_v5_golang/wsdata/protocol"
)

// 订单参数: 下单, 批量下单
type PostOrder struct {
	// 产品 ID
	InstId string `json:"instId"`
	// 交易模式
	TdMode string `json:"tdMode"`
	// 保证金币种
	Ccy string `json:"ccy"`
	// 用户提供的订单 ID
	ClOrdId string `json:"clOrdId"`
	// 订单标签
	Tag string `json:"tag"`
	// 订单方向
	Side string `json:"side"`
	// 持仓方向
	PosSide string `json:"posSide"`
	// 订单类型
	OrdType string `json:"ordType"`
	// 买入或卖出的数量
	Sz string `json:"sz"`
	// 委托价
	Px string `json:"px"`
	// 是否只减仓
	ReduceOnly bool `json:"reduceOnly"`
	// 市价单委托数量的类型
	TgtCcy string `json:"tgtCcy"`
}

// 订单请求参数: 下单, 批量下单
type PostOrderMessage struct {
	// 消息的唯一标识
	Id string `json:"id"`
	// 支持的业务操作
	Op string `json:"op"`
	// 请求参数
	Args []PostOrder `json:"args"`
}

// 单个订单参数
func (c *OkxClient) PostSingleOrder(instId, tdMode, clOrdId, side, posSide, ordType, sz, px string) PostOrder {
	// 登录参数设置
	singleOrder := &PostOrder{
		// 产品 ID
		InstId: instId,
		// 交易模式
		TdMode: tdMode,
		// 用户提供的订单 ID
		ClOrdId: clOrdId,
		// 订单方向
		Side: side,
		// 持仓方向
		PosSide: posSide,
		// 订单类型
		OrdType: ordType,
		// 买入或卖出的数量
		Sz: sz,
		// 委托价
		Px: px,
	}
	// 返回参数
	return *singleOrder
}

// 批量下单
func (c *OkxClient) PostOrders(id, op string, args []PostOrder, dr *DataRepo) error {
	// orders 上锁
	dr.Mu.Lock()
	// orders 添加指针, 指向特殊结构体, 等待交易所返回数据后解除特殊结构体
	for i := 0; i < len(args); i++ {
		// 判断 PingPong 订单类型
		if args[i].InstId == config.PingPongInstID {
			// 跳过操作
			continue
		}
		// 创建特殊订单
		newOrder := &Orders{
			// 订单号
			ClOrdId: args[i].ClOrdId,
			// 订单状态: 本地
			State: "local",
		}
		// 添加到订单列表
		dr.OrdersData[args[i].ClOrdId] = newOrder
	}
	// orders 解锁
	defer dr.Mu.Unlock()

	// 获取参数
	batchOrders := &PostOrderMessage{
		// 唯一标识
		Id: id,
		// 业务操作
		Op: op,
		// 请求参数
		Args: args,
	}
	// 错误提示
	var err error
	// 发起交易
	err = c.send(batchOrders)
	// 错误提示
	if err != nil {
		// 错误提示
		log.Fatalf("[错误提示] 挂单请求失败: %v", err)
		// 返回
		return err
	}
	// 成功提示
	log.Printf("[成功提示] 挂单请求成功")
	// 返回
	return nil
}

// 撤单参数: 撤单, 批量撤单
type CancelOrder struct {
	// 产品 ID
	InstId string `json:"instId"`
	// 订单 ID
	OrdId string `json:"ordId"`
	// 用户提供的订单 ID
	ClOrdId string `json:"clOrdId"`
}

// 撤单请求参数: 撤单, 批量撤单
type CancelOrderMessage struct {
	// 消息的唯一标识
	Id string `json:"id"`
	// 支持的业务操作
	Op string `json:"op"`
	// 请求参数
	Args []CancelOrder `json:"args"`
}

// 撤销订单参数
func (c *OkxClient) CancelSingleOrder(instId, OrdId, clOrdId string) CancelOrder {
	// 登录参数设置
	cancelSingleOrder := &CancelOrder{
		// 产品 ID
		InstId: instId,
		// 订单 ID
		OrdId: OrdId,
		// 用户提供的订单 ID
		ClOrdId: clOrdId,
	}
	// 返回参数
	return *cancelSingleOrder
}

// 批量撤单
func (c *OkxClient) CancelOrders(id, op string, args []CancelOrder, dr *DataRepo) error {
	// orders 上锁
	dr.Mu.Lock()
	// 逐个订单判定
	for i := 0; i < len(args); i++ {
		// orders 判断 Key 值是否存在
		if val, ok := dr.OrdersData[args[i].ClOrdId]; ok {
			// 若存在, orders 判断 Key 值是否仍然指向特殊结构体, 若是, 说明订单未成功发出, 删除 Key
			if val.State == "local" {
				// 说明订单未成功发出, 删除 Key
				delete(dr.OrdersData, args[i].ClOrdId)
			}
			// 若存在, 且指向非空地址, 则无需操作, 因为此时撤销订单后, 会自动删除指针
		}
		// 若不存在, 说明已经撤销或成交, 则无需操作
	}
	// orders 解锁
	defer dr.Mu.Unlock()

	// 获取参数
	batchOrders := &CancelOrderMessage{
		// 唯一标识
		Id: id,
		// 业务操作
		Op: op,
		// 请求参数
		Args: args,
	}
	// 错误提示
	var err error
	// 发起交易
	err = c.send(batchOrders)
	// 错误提示
	if err != nil {
		// 错误提示
		log.Fatalf("[错误提示] 撤单请求失败: %v", err)
		// 返回
		return err
	}
	// 成功提示
	log.Printf("[成功提示] 撤单请求成功")
	// 返回
	return nil
}
