# Author: Wiger
# Date: 2022-04-05
# Description: 分批读取交易所提供的 trade 数据,
#              提取其中盘口 5 档数据, 转存至多个子文件

# 导入依赖
import os
import pandas as pd

# 目标交易对
swap_name = 'DOT-USDT-SWAP'
# 文件夹路径
csv_folder = '/Users/wiger/Documents/historydata/'
# 文件名称
file_path = 'allswap-trades-2022-02-09.csv'
# 文件夹名称 例如: 'allfuture-trades-2022-02-10.csv'
file_pre_path = file_path[:-4]
# 赋予列名
cols = ['instID', 'unknown', 'side', 'size', 'price', 'timestamp']

# 单次读取行数
chunk_patch = 100000
# 分批读取
chunker = pd.read_csv(csv_folder + file_path, chunksize=chunk_patch, header=None, names=cols)
# 批次数目
patch_num = 1
# 分批读取
for item in chunker:
    # 提取部分行
    part_item = item.loc[item['instID'] == swap_name]
    # 文件夹路径
    folder_path = csv_folder + file_pre_path + '/'
    # 判断文件夹是否存在
    if not os.path.exists(folder_path):
        # 创建文件夹
        os.makedirs(folder_path)
    # 分批文件名
    file_name = csv_folder + file_pre_path + '/' + file_pre_path + '.' + str(patch_num) + '.csv'
    # 保存文件
    part_item.to_csv(file_name, header=1, index=0)
    # 打印进度
    print('已提取文件: {0}'.format(file_name))
    # 更新批次数
    patch_num += 1
