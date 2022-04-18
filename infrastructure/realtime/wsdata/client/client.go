package client

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wiger123/okex_v5_golang/config"
	. "github.com/wiger123/okex_v5_golang/utils"
	. "github.com/wiger123/okex_v5_golang/wsdata/protocol"
)

// 频道订阅参数
type Arg struct {
	// 频道名
	Channel string `json:"channel"`
	// 产品类型
	InstType string `json:"instType"`
	// 标的指数
	Uly string `json:"uly"`
	// 产品 ID
	InstID string `json:"instId"`
}

// 信息处理
type MessageHandler func(m PushMessage)

// 客户端
type OkxClient struct {
	// 并发锁
	mux sync.RWMutex
	// websocket 连接
	conn *websocket.Conn
	// 频道订阅参数
	channels []Arg
	// 信息处理字典
	handlers map[string]MessageHandler
}

// 账户信息
type LoginArg struct {
	// API Key
	APIKey string `json:"apiKey"`
	// API Passphrase
	Passphrase string `json:"passphrase"`
	// 时间戳
	Timestamp string `json:"timestamp"`
	// 签名字符串
	Sign string `json:"sign"`
}

// 登录请求
type LoginRequest struct {
	// 操作
	Op string `json:"op"`
	// 账户列表
	Args []LoginArg `json:"args"`
}

// 订阅请求
type SubscribeRequest struct {
	// 操作
	Op string `json:"op"`
	// 订阅频道列表
	Args []Arg `json:"args"`
}

// Websocket 信息内容
type MessageProfile struct {
	// 参数
	Arg struct {
		// 频道名
		Channel string `json:"channel"`
		// 产品 ID
		InstID string `json:"instId"`
	} `json:"arg"`
}

// 创建新的客户端
func NewOkxClient(url string) (*OkxClient, error) {
	// 发起 websocket 连接
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	// 报错
	if err != nil {
		// 错误提示
		log.Fatalf("[错误提示] OKX 客户端登录失败: %v", err)
		// 返回错误
		return nil, err
	}
	// 成功提示
	log.Printf("[成功提示] OKX 客户端登录成功")
	// 成功返回客户端
	return &OkxClient{
		// 连接通道
		conn: conn,
		// 信息处理
		handlers: make(map[string]MessageHandler),
	}, nil
}

// 登陆
func (c *OkxClient) Login() error {
	// 时间戳
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// request 路径
	message := PreHashString(timestamp, "GET", "/users/self/verify", "")
	// HMAC SHA256
	sign, err := HmacSha256Base64Signer(message, config.SecretKey)
	// 错误提示
	if err != nil {
		// 返回
		return err
	}
	// 登录参数设置
	request := &LoginRequest{
		// 操作
		Op: "login",
		// 账户列表
		Args: []LoginArg{
			{
				// API Key
				APIKey: config.ApiKey,
				// API Passphrase
				Passphrase: config.PassPhrase,
				// 时间戳
				Timestamp: timestamp,
				// 签名字符串
				Sign: sign,
			},
		},
	}
	// 发送登录请求
	err = c.send(request)
	// 错误提示
	if err != nil {
		// 错误提示
		log.Fatalf("[错误提示] OKX 私有频道登陆失败: %v", err)
		// 返回
		return err
	}
	// 成功提示
	log.Printf("[成功提示] OKX 私有频道登录成功")
	// 返回
	return nil
}

// 订阅频道参数初始化
func (c *OkxClient) Subscribe(channel, instType, uly, instID string, handler MessageHandler) {
	// 订阅频道参数
	c.channels = append(c.channels, Arg{
		// 频道名
		Channel: channel,
		// 产品类型
		InstType: instType,
		// 标的指数
		Uly: uly,
		// 产品 ID
		InstID: instID,
	})
	// 信息处理字典赋值
	c.handlers[c.channelKey(channel, instID)] = handler
}

// 读取 Websocket 数据并处理
func (c *OkxClient) ReadWebsocket() {
	// 读取 websocket 数据
	_, data, err := c.conn.ReadMessage()

	// 读取信息失败
	if err != nil {
		// 错误提示
		log.Printf("[错误提示] Websocket 数据读取失败: %v", err)
	}

	// 初始化解析数据
	var profile MessageProfile
	// 解析数据
	if err := json.Unmarshal(data, &profile); err != nil {
		// 错误提示
		log.Printf("[错误提示] Websocket 数据解析失败: %v", err)
	}

	// 初始化数据处理接口
	var message PushMessage
	// 数据分类处理
	switch profile.Arg.Channel {
	// 交易数据
	case "trades":
		// 交易数据初始化
		var tm TradeMessage
		// 数据解析
		json.Unmarshal(data, &tm)
		// 数据内容
		message = &tm
	// 盘口数据
	case "books5":
		// 盘口数据初始化
		var bm Book5Message
		// 数据解析
		json.Unmarshal(data, &bm)
		// 数据内容
		message = &bm
	// 账户数据
	case "account":
		// 账户数据初始化
		var am AccountMessage
		// 数据解析
		json.Unmarshal(data, &am)
		// 数据内容
		message = &am
	// 持仓数据
	case "positions":
		// 持仓数据初始化
		var pm PositionsMessage
		// 数据解析
		json.Unmarshal(data, &pm)
		// 数据内容
		message = &pm
	// 订单数据
	case "orders":
		// 订单数据初始化
		var om OrdersMessage
		// 数据解析
		json.Unmarshal(data, &om)
		// 数据内容
		message = &om
	// 其他情况
	default:
		// 普通提示
		// log.Printf("[普通提示] Websocket 请求响应: %v", string(data))
		// 返回
		return
	}

	// 提取频道名称和产品 ID
	channel, instID := message.ChannelAndInstID()
	// 频道名称和产品 ID 字符串格式化
	channelKey := c.channelKey(channel, instID)
	// 选取指定信息处理器
	handler := c.handlers[channelKey]
	// 未找到信息处理器
	if handler == nil {
		// 普通提示
		log.Printf("[普通提示] 未知信息处理器: %v", channelKey)
		// 返回
		return
	}

	// 处理信息
	handler(message)
}

// 循环读取 Websocket 数据
func (c *OkxClient) ReadWebsocketLoop() {
	// 循环
	for {
		c.ReadWebsocket()
	}
}

// 发起 Websocket 连接, 解析数据, 并执行数据存储
func (c *OkxClient) Run() error {
	// 订阅 channel
	if err := c.send(&SubscribeRequest{
		// 操作
		Op: "subscribe",
		// 订阅频道列表
		Args: c.channels,
	}); err != nil {
		// 错误提示
		log.Fatalf("[错误提示] Websocket 订阅失败: %v", err)
		// 返回错误
		return err
	}
	// 成功提示
	log.Printf("[成功提示] Websocket 订阅成功")

	// 返回
	return nil
}

// 发送请求
func (c *OkxClient) send(message interface{}) error {
	// 显示
	// log.Printf("[普通提示] 数据解析内容: %v", message)
	// 结构体转为 json 格式
	data, err := json.Marshal(message)
	// 解析错误
	if err != nil {
		// 错误提示
		log.Fatalf("[错误提示] 频道参数转为 json 格式失败: %v", err)
		// 返回错误
		return err
	}

	// 上锁
	c.mux.Lock()
	// 函数结束前解锁
	defer c.mux.Unlock()
	// 发送请求
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// 频道 Key 值格式化
func (c *OkxClient) channelKey(channel, instID string) string {
	// 字符串格式化
	return fmt.Sprintf("channel:%v, instID:%v", channel, instID)
}

// 关闭客户端
func (c *OkxClient) Shutdown() {
	// 关闭连接
	c.conn.Close()
	// 成功提示
	log.Printf("[成功提示] OKX 客户端连接关闭")
}
