# Author: Wiger
# Date: 2022-03-19
# Description: 分批读取交易所提供的 50 档 ticker 数据,
#              提取其中盘口 5 档数据, 转存至多个子文件

# 导入依赖
import os
import pandas as pd

# 文件夹路径
csv_folder = 'D:/okx_ticker_data/dot/'
# 文件名称
file_path = '20220209DOT-USDT-SWAP.OK.csv'
# 文件夹名称 例如: '20220216DOT-USDT-SWAP.OK'
file_pre_path = file_path[:-4]
# 列名
columns_path = 'columns.txt'
# 提取列
columns_sel = ['timeMs', 'exchTimeMs', 'maxLevel', 
               'askPx1', 'askCnt1', 'askSz1', 
               'askPx2', 'askCnt2', 'askSz2', 
               'askPx3', 'askCnt3', 'askSz3', 
               'askPx4', 'askCnt4', 'askSz4', 
               'askPx5', 'askCnt5', 'askSz5',
               'bidPx1', 'bidCnt1', 'bidSz1', 
               'bidPx2', 'bidCnt2', 'bidSz2', 
               'bidPx3', 'bidCnt3', 'bidSz3', 
               'bidPx4', 'bidCnt4', 'bidSz4', 
               'bidPx5', 'bidCnt5', 'bidSz5']

# 单次读取行数
chunk_patch = 100000
# 分批读取
chunker = pd.read_csv(csv_folder + file_path, chunksize=chunk_patch)
# 批次数目
patch_num = 1
# 分批读取
for item in chunker:
    # 首次读取需要读取列名
    # 之后只需要提取所需特定列
    # 打开文件
    # with open(csv_folder + columns_path, 'w') as f:
        # 记录列名
        # f.write(str(list(item.columns)))
    # 关闭文件
    # f.close()

    # 提取部分列
    part_item = item[columns_sel]
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
