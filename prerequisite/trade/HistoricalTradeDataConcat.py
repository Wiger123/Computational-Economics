# Author: Wiger
# Date: 2022-04-05
# Description: 将多个 trade 数据合并为一个
#              并只保留一个表头

# 导入依赖
import shutil
import glob

# 输入文件夹路径
in_path = 'D:/okx_ticker_data/tradeData/allfuture-trades-2022-02-10/'
# 输出文件夹路径
out_path = 'D:/okx_ticker_data/dot_trade/'
# 输出文件名称
out_name = '20220210DOT-USDT-220325.OK.trade.csv'
# 路径下文件
all_files = glob.glob(in_path + "/*.csv")
# 按照文件名称排序
all_files.sort()
# 输出文件路径
with open(out_path + out_name, 'wb') as out_file:
    # 迭代加入
    for i, fname in enumerate(all_files):
        # 打开文件
        with open(fname, 'rb') as in_file:
            # 除了第一个表格读取表头
            if i != 0:
                # 后续表格跳过表头
                in_file.readline()
            # 在不解析的情况下, 将文件的其余内容, 从输入复制到输出
            shutil.copyfileobj(in_file, out_file)
            # 打印
            print('已合并文件数: {0}'.format(i + 1))
        # 关闭文件
        in_file.close()
# 关闭文件
out_file.close()
