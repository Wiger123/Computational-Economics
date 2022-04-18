package strategy

import (
	"log"
	"time"

	. "github.com/wiger123/okex_v5_golang/config"
	. "github.com/wiger123/okex_v5_golang/database"
	. "github.com/wiger123/okex_v5_golang/utils"
	. "github.com/wiger123/okex_v5_golang/wsdata/client"
	. "github.com/wiger123/okex_v5_golang/wsdata/protocol"
)

// 印钞机策略
func PrintMoney(c *OkxClient, dataRepo *DataRepo) {
	// 加载数据库
	DatabaseLoader(dataRepo)
	// 构建数据库
	printMoneyData := NewPrintMoney()
	// 策略循环
	for {
		// 核心策略
		PrintMoneyCore(*printMoneyData, dataRepo)
		// 等待
		time.Sleep(2 * time.Second)
	}
}

// 数据库加载
func DatabaseLoader(dataRepo *DataRepo) {
	// 数据收集中, 等待启动
	for {
		// 判断交易数据数目 盘口数据数目
		if len(dataRepo.TradeData) < Ntrade || len(dataRepo.Book5Data) < NBook5s {
			// 提示
			log.Printf("[普通提示] 基础数据正在收集中, 策略即将启动, 请等待: Trade: %v / %v, Book5: %v / %v", len(dataRepo.TradeData), Ntrade, len(dataRepo.Book5Data), NBook5s)
		} else {
			// 跳出循环
			break
		}
		// 等待
		time.Sleep(500 * time.Millisecond)
	}
	// 策略启动
	log.Printf("[成功提示] 策略启动")
}

// 加权交易量, 获取交易量时间
func WeightVol(lastVol float64, lastTradeTime int64, dataRepo *DataRepo) (float64, int64) {
	// 最新时间
	var newTradeTime int64
	// 初始化
	newTradeTime = lastTradeTime
	// 当前交易记录求和
	var sumVol float64
	// 初始化
	sumVol = 0
	// 循环
	for i := 0; i < Ntrade; i++ {
		// 时间戳
		var time = String2Int64(dataRepo.TradeData[i].Ts)
		// 更新最新时间
		newTradeTime = MaxInt64(newTradeTime, time)
		// 交易时间大于上次交易时间
		if time > lastTradeTime {
			// 交易量累加
			sumVol += String2Float64(dataRepo.TradeData[i].Sz)
		}
	}
	// 加权求和
	var newVol = 0.7*lastVol + 0.3*sumVol
	// 最新交易量, 最新交易时间
	log.Printf("[普通提示] 交易量: %.3f, 交易时间: %v", newVol, newTradeTime)
	// 返回数据
	return newVol, newTradeTime
}

// 获取仓位数据, 平衡仓位
func BalanceAccount(dataRepo *DataRepo) float64 {
	// 仓位价值
	tokenValue := dataRepo.TokenAmt * String2Float64(dataRepo.Book5Data[0].Bids[0][0])
	// 仓位比例
	res := tokenValue / (tokenValue + dataRepo.UsdtAmt)
	// 仓位小于平衡
	if res < BalancePos-BalanceRel {
		// 挂小买单: Price: Bids[0] + 0.000 / 0.001 / 0.002  Size: 0.01
		// 定时撤单
	}
	// 仓位大于平衡
	if res > BalancePos+BalanceRel {
		// 挂小卖单: Price: Asks[0] - 0.000 / 0.001 / 0.002  Size: 0.01
		// 定时撤单
	}
	// 仓位
	log.Printf("[普通提示] Token: %v, Token 余额: %v, USDT: %v, USDT 余额: %v, 仓位占比: %v", TokenInstID, dataRepo.TokenAmt, UsdtInstID, dataRepo.UsdtAmt, res)
	// 返回
	return res
}

// 核心策略
func PrintMoneyCore(printMoneyData PrintMoneyAttr, dataRepo *DataRepo) {
	// 若当前有挂单
	if len(dataRepo.OrdersData) > 0 {
		// 返回
		return
	}
	// 计数器
	printMoneyData.NumTick++
	// 获取加权交易量, 最近交易时间更新
	printMoneyData.Vol, printMoneyData.LastTradeTime = WeightVol(printMoneyData.Vol, printMoneyData.LastTradeTime, dataRepo)
	// 平衡仓位
	printMoneyData.P = BalanceAccount(dataRepo)
	// 爆发价格
	var burstPrice = printMoneyData.Prices[NBook5sAvg-1] * BurstThresholdPct
	// 牛市变量
	var bull = false
	// 熊市变量
	var bear = false
	// 交易数量
	var tradeAmount float64
	// 显示数据
	log.Printf("[普通提示] 爆发价格: %v, 牛市变量: %v, 熊市变量: %v, 交易数量: %v", burstPrice, bull, bear, tradeAmount)
	// 价格倒数 1
	newPrice1 := printMoneyData.Prices[NBook5sAvg-1]
	// 价格倒数 2
	newPrice2 := printMoneyData.Prices[NBook5sAvg-2]
	// 价格倒数 6 - 1 位最大值
	maxLast6to1 := MaxSlice(printMoneyData.Prices[NBook5sAvg-6 : NBook5sAvg-1])
	// 价格倒数 6 - 1 位最小值
	minLast6to1 := MinSlice(printMoneyData.Prices[NBook5sAvg-6 : NBook5sAvg-1])
	// 价格倒数 6 - 2 位最大值
	maxLast6to2 := MaxSlice(printMoneyData.Prices[NBook5sAvg-6 : NBook5sAvg-2])
	// 价格倒数 6 - 2 位最小值
	minLast6to2 := MinSlice(printMoneyData.Prices[NBook5sAvg-6 : NBook5sAvg-2])
	// 判断牛熊
	if printMoneyData.NumTick > 2 &&
		(newPrice1-maxLast6to1 > burstPrice ||
			newPrice1-maxLast6to2 > burstPrice &&
				newPrice1 > newPrice2) {
		// 牛市变量
		bull = true
		// 交易数量
		tradeAmount = dataRepo.UsdtAmt / printMoneyData.BidPrice * 0.99
	} else if printMoneyData.NumTick > 2 &&
		(newPrice1-minLast6to1 < -burstPrice ||
			newPrice1-minLast6to2 < -burstPrice &&
				newPrice1 < newPrice2) {
		// 熊市变量
		bear = true
		// 交易数量
		tradeAmount = dataRepo.TokenAmt
	}

	// 缩减交易量: 历史交易量未达阈值
	if printMoneyData.Vol < BurstThresholdVol {
		// 交易量
		tradeAmount *= printMoneyData.Vol / BurstThresholdVol
	}
	// 缩减交易量: 循环次数未达阈值
	if printMoneyData.NumTick < 5 {
		// 交易量
		tradeAmount *= 0.8
	}
	// 缩减交易量: 循环次数未达阈值
	if printMoneyData.NumTick < 10 {
		// 交易量
		tradeAmount *= 0.8
	}
	// 非牛市熊市
	if (!bull && !bear) || tradeAmount < MinStock {
		// 返回
		return
	}

	// 交易价格
	var tradePrice float64
	// 根据牛市熊市确定价格
	if bull == true {
		// 牛市交易价格
		tradePrice = printMoneyData.BidPrice
	} else if bear == true {
		// 熊市交易价格
		tradePrice = printMoneyData.AskPrice
	} else {
		// 其他情况返回
		return
	}
	// 价格
	log.Panicf("[普通提示] 牛市: %v, 熊市: %v, 交易价格: %v", bull, bear, tradePrice)

	// 发起交易
	if bull == true {
		// 牛市买入
		// 定时撤单
	} else if bear == true {
		// 熊市卖出
		// 定时撤单
	} else {
		// 其他情况返回
		return
	}

	// 更新计数
	printMoneyData.NumTick = 0
}
