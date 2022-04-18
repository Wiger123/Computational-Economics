package strategy

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/wiger123/okex_v5_golang/config"
	. "github.com/wiger123/okex_v5_golang/database"
	. "github.com/wiger123/okex_v5_golang/utils"
	. "github.com/wiger123/okex_v5_golang/wsdata/client"
)

/**
	策略逻辑:
		多因子线性组合判断短期趋势,
		趋势为涨, 则在盘口挂多单, 平空仓,
	   	趋势为跌, 则在盘口挂空单, 平多仓,
	   	仓位量由交易量, 投入总金额, 最低金额比例决定

	因子1:
		1. 对最近时间和最远时间做权重划分: [0.2, 0.4, 0.6, 0.8, 1]
		2. 对时间窗口内交易量进行加权统计
		3. 阈值设定, 根据结果与阈值的比较, 判定是否操作
**/

// 执行策略
func OnStrategy1(c *OkxClient, dataRepo *DataRepo) {
	// 循环
	for {
		// 判断交易数据数目 盘口数据数目
		if len(dataRepo.TradeData) < config.Ntrade || len(dataRepo.Book5Data) < config.NBook5s {
			// 提示
			log.Printf("[普通提示] 基础数据正在收集中, 策略即将启动, 请等待: Trade: %v / %v, Book5: %v / %v", len(dataRepo.TradeData), config.Ntrade, len(dataRepo.Book5Data), config.NBook5s)
		} else {
			// 跳出循环
			break
		}
		// 等待
		time.Sleep(500 * time.Millisecond)
	}

	// 策略启动
	log.Printf("[成功提示] 策略启动")

	// 构建权重参数
	var weightList []float64
	// 间距
	var weightInterval = (config.MaxWeight - config.MinWeight) / (config.NumLevel - 1)
	// 添加元素
	for i := 0; i < config.NumLevel; i++ {
		// 计算结果
		var res = config.MinWeight + weightInterval*float64(i)
		// 保留小数
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", res), 64)
		// 添加元素
		weightList = append(weightList, res)
	}
	// 显示
	log.Printf("[成功提示] 权重分配: %v", weightList)

	// 构建权重参数
	var postList []float64
	// 间距
	var postInterval = (config.MaxPost - config.MinPost) / (config.NumPost - 1)
	// 添加元素
	for i := 0; i < config.NumPost; i++ {
		// 计算结果
		var res = config.MinPost + postInterval*float64(i)
		// 保留小数
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", res), 64)
		// 添加元素
		postList = append(postList, res)
	}
	// 显示
	log.Printf("[成功提示] 仓位分配: %v", postList)

	// 数据间隔
	var dataInterval = config.Ntrade / config.NumLevel
	// 显示
	log.Printf("[成功提示] 数据档位数: %v  每档数据量: %v", config.NumLevel, dataInterval)

	// 循环
	for {
		// sell 权重
		var sellWeight float64
		// buy 权重
		var buyWeight float64
		// 计算买卖双方动向
		for i := config.Ntrade - 1; i >= 0; i-- {
			// 判断方向
			if dataRepo.TradeData[i].Side == "buy" {
				// 量
				var perSize, _ = strconv.ParseFloat(dataRepo.TradeData[i].Sz, 64)
				// 档位
				var perLevel = int(math.Floor(float64(i) / float64(dataInterval)))
				// 加权交易量
				var perWeightSize = weightList[perLevel] * perSize
				// 累计
				buyWeight += perWeightSize
			} else {
				// 量
				var perSize, _ = strconv.ParseFloat(dataRepo.TradeData[i].Sz, 64)
				// 档位
				var perLevel = int(math.Floor(float64(i) / float64(dataInterval)))
				// 加权交易量
				var perWeightSize = weightList[perLevel] * perSize
				// 累计
				sellWeight += perWeightSize
			}
		}
		// 买卖量
		// log.Printf("[成功提示] 买单加权量: %v  卖单加权量: %v", buyWeight, sellWeight)

		// 挂多 平空
		if buyWeight-config.Ratio*sellWeight > 0 && buyWeight > config.MinTradeVolume {
			// 订单聚合
			var orders []PostOrder
			// 订单 ID
			var cltId1 = GetRandString(config.ClOrdIdLength)
			// 订单 ID
			var cltId2 = GetRandString(config.ClOrdIdLength)
			// 有空仓
			if dataRepo.PositionsShortData.AvailPos != "" {
				// 数量
				var coverSize = dataRepo.PositionsShortData.AvailPos
				// 价格
				var coverPrice = dataRepo.Book5Data[config.NBook5s-1].Bids[config.CoverShortLevel][0]
				// 平仓
				var order1 = c.PostSingleOrder(config.InstID, config.TdMode, cltId1, "buy", "short", "post_only", coverSize, coverPrice)
				// 添加订单
				orders = append(orders, order1)
			}
			// 数量
			var postSize = strconv.FormatFloat(postList[Min(int(math.Floor(sellWeight*10/config.MaxRef)), len(postList)-1)], 'f', config.FloatPrec, 64)
			// 价格
			var postPrice = dataRepo.Book5Data[config.NBook5s-1].Bids[config.BidsLevel][0]
			// 开仓
			var order2 = c.PostSingleOrder(config.InstID, config.TdMode, cltId2, "buy", "long", "post_only", postSize, postPrice)
			// 添加订单
			orders = append(orders, order2)
			// 批量下单
			c.PostOrders("posttest", "batch-orders", orders, dataRepo)
			// 显示
			// log.Printf("[成功提示] 平空  挂多: %v", postSize)
			// 显示
			// log.Printf("[成功提示] 挂单价格: %v", postPrice)

			// 撤单定时
			durationOfTime := time.Duration(config.TimeCancel) * time.Millisecond
			// 取消订单函数
			f := func() {
				// 订单聚合
				var corders []CancelOrder
				// 撤销订单
				var corder1 = c.CancelSingleOrder(config.InstID, "", cltId1)
				// 撤销订单
				var corder2 = c.CancelSingleOrder(config.InstID, "", cltId2)
				// 添加订单
				corders = append(corders, corder1)
				// 添加订单
				corders = append(corders, corder2)
				// 批量撤单
				c.CancelOrders("posttest", "batch-cancel-orders", corders, dataRepo)
			}
			// 计时器
			time.AfterFunc(durationOfTime, f)
		}

		// 挂空 平多
		if sellWeight-config.Ratio*buyWeight > 0 && sellWeight > config.MinTradeVolume {
			// 订单聚合
			var orders []PostOrder
			// 订单 ID
			var cltId1 = GetRandString(config.ClOrdIdLength)
			// 订单 ID
			var cltId2 = GetRandString(config.ClOrdIdLength)
			// 有多仓
			if dataRepo.PositionsLongData.AvailPos != "" {
				// 数量
				var coverSize = dataRepo.PositionsLongData.AvailPos
				// 价格
				var coverPrice = dataRepo.Book5Data[config.NBook5s-1].Asks[config.CoverLongLevel][0]
				// 平仓
				var order1 = c.PostSingleOrder(config.InstID, config.TdMode, cltId1, "sell", "long", "post_only", coverSize, coverPrice)
				// 添加订单
				orders = append(orders, order1)
			}
			// 数量
			var postSize = strconv.FormatFloat(postList[Min(int(math.Floor(sellWeight*10/config.MaxRef)), len(postList)-1)], 'f', config.FloatPrec, 64)
			// 价格
			var postPrice = dataRepo.Book5Data[config.NBook5s-1].Asks[config.AsksLevel][0]
			// 开仓
			var order2 = c.PostSingleOrder(config.InstID, config.TdMode, cltId2, "sell", "short", "post_only", postSize, postPrice)
			// 添加订单
			orders = append(orders, order2)
			// 批量下单
			c.PostOrders("posttest", "batch-orders", orders, dataRepo)
			// 显示
			// log.Printf("[成功提示] 平多  挂空: %v", postSize)
			// 显示
			// log.Printf("[成功提示] 挂单价格: %v", postPrice)

			// 撤单定时
			durationOfTime := time.Duration(config.TimeCancel) * time.Millisecond
			// 取消订单函数
			f := func() {
				// 订单聚合
				var corders []CancelOrder
				// 撤销订单
				var corder1 = c.CancelSingleOrder(config.InstID, "", cltId1)
				// 撤销订单
				var corder2 = c.CancelSingleOrder(config.InstID, "", cltId2)
				// 添加订单
				corders = append(corders, corder1)
				// 添加订单
				corders = append(corders, corder2)
				// 批量撤单
				c.CancelOrders("posttest", "batch-cancel-orders", corders, dataRepo)
			}
			// 计时器
			time.AfterFunc(durationOfTime, f)
		}

		// 仓位信息
		// log.Printf("[成功提示] 多仓信息: %v  空仓信息: %v", dataRepo.PositionsLongData, dataRepo.PositionsShortData)

		// 挂单信息
		// log.Printf("[成功提示] 挂单信息: %v", dataRepo.OrdersData)

		// 等待
		time.Sleep(100 * time.Millisecond)
	}
}
