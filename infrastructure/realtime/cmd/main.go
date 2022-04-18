package main

import (
	"time"

	"github.com/wiger123/okex_v5_golang/config"
	. "github.com/wiger123/okex_v5_golang/database"
	. "github.com/wiger123/okex_v5_golang/strategy"
	. "github.com/wiger123/okex_v5_golang/wsdata/client"
)

// 主函数
func main() {
	// 数据库初始化
	dataRepo := NewDataRepo()

	// 创建 okx 客户端: 公共频道
	publicClient, _ := NewOkxClient(config.PublicURL)
	// 创建 okx 客户端: 私有频道
	privateClient, _ := NewOkxClient(config.PrivateURL)
	// 私有频道登陆
	privateClient.Login()

	// 等待登陆成功
	time.Sleep(5 * time.Second)

	// 公共频道添加订阅
	// 交易频道
	publicClient.Subscribe("trades", "", "", config.InstID, dataRepo.HandleMessage)
	// 盘口频道
	publicClient.Subscribe("books5", "", "", config.InstID, dataRepo.HandleMessage)
	// 私有频道添加订阅
	// 账户频道
	privateClient.Subscribe("account", "", "", "", dataRepo.HandleMessage)
	// 持仓频道
	privateClient.Subscribe("positions", config.InstType, "", config.InstID, dataRepo.HandleMessage)
	// 订单频道
	privateClient.Subscribe("orders", config.InstType, "", config.InstID, dataRepo.HandleMessage)

	// 公共频道订阅
	publicClient.Run()
	// 私有频道订阅
	privateClient.Run()

	// 公共频道数据解析并处理
	go publicClient.ReadWebsocketLoop()
	// 私有频道数据解析并处理
	go privateClient.ReadWebsocketLoop()

	// 等待
	time.Sleep(3 * time.Second)

	// 私有频道保持连接
	go PingPong(privateClient, dataRepo)
	// 执行策略
	go PrintMoney(privateClient, dataRepo)

	// 等待
	time.Sleep(19990726 * time.Second)

	// 关闭公共频道客户端
	publicClient.Shutdown()
	// 关闭私有频道客户端
	privateClient.Shutdown()
}
