# Computational-Economics
- Computational Economics: High Frequency Trading
- 计算经济学: 高频交易

### 项目目录
- prerequisite: 预处理文件, 包含 ticker, trade 等数据的预处理操作
  - ticker: 盘口数据预处理
  - trade: 交易数据预处理
- infrastructure: 框架选型, 包含回测框架, 实盘框架, 主要为使用文档, 多以文字调研结果呈现
  - backtest: 回测框架
    - backup: 回测模块
  - realtime: 实盘框架
    - okgo: golang 编写的 okex 交易框架

### 因子
- 因子: 数据 => 构造指标 => 根据指标的数值, 决定是否交易

### 传统因子
- 可以尝试复现: 马丁格尔, MACD, 布林带, 波动率

### 批量构造
- +-*/
- dydx
- 偏度, 峰度
- 傅里叶
- ...

### 数据资源
- 拷贝获取

### 团队文档
- 石墨文档: https://shimo.im/tables/wV3VVj0jprFnzj3y?table=P6eddwEL9gW&view=g1MkIQ4l9XB#/
- 任务分工及进度:
![image](https://user-images.githubusercontent.com/31722033/161755000-20d72790-2db6-4aa1-b608-ee4cdb2736ce.png)
- 2022.04.20:
![微信截图_20220420103512](https://user-images.githubusercontent.com/31722033/164135679-bc83ddfb-2622-4b1f-a9fd-6c0b3f173532.png)
