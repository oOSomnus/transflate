#!/bin/bash

# 检查是否安装了 cloc
if ! command -v cloc &> /dev/null; then
  echo "cloc 未安装，请先安装 cloc 后重试。"
  exit 1
fi

# 确保脚本在仓库根目录运行
if [ ! -d ".git" ]; then
  echo "请确保在 Git 仓库的根目录运行该脚本。"
  exit 1
fi

echo "正在统计代码仓库：$(pwd)"

# 使用 cloc 统计代码信息
cloc_output=$(cloc .)

# 输出统计结果
echo "统计结果："
echo "$cloc_output"
