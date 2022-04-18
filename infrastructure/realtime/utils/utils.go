package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 取小
func Min(a, b int) int {
	// a 比 b 小
	if a < b {
		// 返回 a
		return a
	}
	// 否则返回 b
	return b
}

// 取大
func Max(a, b int) int {
	// a 比 b 大
	if a > b {
		// 返回 a
		return a
	}
	// 否则返回 b
	return b
}

// 取小
func MinInt64(a, b int64) int64 {
	// a 比 b 小
	if a < b {
		// 返回 a
		return a
	}
	// 否则返回 b
	return b
}

// 取大
func MaxInt64(a, b int64) int64 {
	// a 比 b 大
	if a > b {
		// 返回 a
		return a
	}
	// 否则返回 b
	return b
}

// 对字符串进行 HMAC SHA256 加密, 再用 Base 64 编码
func HmacSha256Base64Signer(message string, secretKey string) (string, error) {
	// hmac 算法
	mac := hmac.New(sha256.New, []byte(secretKey))
	// 写入信息
	_, err := mac.Write([]byte(message))
	// 发送错误
	if err != nil {
		// 错误提示
		log.Fatalf("[错误提示] HMAC SHA256 加密失败: %v", err)
		// 返回
		return "", err
	}
	// 返回加密签名
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// 字符串拼接
func PreHashString(timestamp string, method string, requestPath string, body string) string {
	// 字符串拼接
	return timestamp + strings.ToUpper(method) + requestPath + body
}

// 字典库
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 随机种子
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// 生成定长随机字符串
func StringWithCharset(length int, charset string) string {
	// 字节数组初始化
	b := make([]byte, length)
	// 逐个选取
	for i := range b {
		// 逐位赋值
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	// 返回字符串结果
	return string(b)
}

// 定长随机字符串
func GetRandString(length int) string {
	// 返回随机字符串
	return StringWithCharset(length, charset)
}

// 收益率计算
func GetProfitRatio(avgPx, nowPx, lever float64, posSide string) float64 {
	// 收益率
	var pr float64
	// 多
	if posSide == "long" {
		// 计算收益率
		pr = (nowPx - avgPx) / avgPx * lever * 100
	}
	// 空
	if posSide == "short" {
		// 计算收益率
		pr = (avgPx - nowPx) / avgPx * lever * 100
	}
	// 返回结果
	return pr
}

// String 转 Float64
func String2Float64(in string) float64 {
	// 结果
	res, _ := strconv.ParseFloat(in, 64)
	// 返回
	return res
}

// String 转 int
func String2Int64(in string) int64 {
	// 结果
	res, _ := strconv.ParseInt(in, 10, 64)
	// 返回
	return res
}

// Float64 转 String
func Float642String(in float64) string {
	// 结果
	res := strconv.FormatFloat(in, 'f', -1, 64)
	// 返回
	return res
}

// Int64 转 String
func Int642String(in int64) string {
	// 结果
	res := strconv.FormatInt(in, 10)
	// 返回
	return res
}

// Slice 最大值
func MaxSlice(list []float64) float64 {
	// 没有元素
	if len(list) == 0 {
		// 返回空值
		return 0
	}
	// 最大值
	max := list[0]
	// 循环
	for _, v := range list {
		// 判断价格
		if v > max {
			// 更新最大值
			max = v
		}
	}
	// 返回结果
	return max
}

// Slice 最小值
func MinSlice(list []float64) float64 {
	// 没有元素
	if len(list) == 0 {
		// 返回空值
		return 0
	}
	// 最大值
	min := list[0]
	// 循环
	for _, v := range list {
		// 判断价格
		if v < min {
			// 更新最大值
			min = v
		}
	}
	// 返回结果
	return min
}
