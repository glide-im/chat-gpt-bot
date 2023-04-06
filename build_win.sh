#!/bin/bash

# 指定要编译的 Golang 项目的目录
PROJECT_DIR=.

# 指定编译后的二进制文件的名称和输出目录
OUTPUT_NAME=bot.exe
OUTPUT_DIR=./

# 编译 Golang 项目
echo "编译 $PROJECT_DIR..."
GOOS=windows GOARCH=amd64 go build -o "$OUTPUT_DIR/$OUTPUT_NAME" "$PROJECT_DIR"

# 检查编译结果
if [ $? -eq 0 ]; then
    echo "编译成功！"
else
    echo "编译失败！"
fi

wait