#!/bin/bash

echo "🚀 Go数据库操作演示启动脚本"
echo "================================"

# 检查MySQL是否运行
echo "📋 检查MySQL服务状态..."
if brew services list | grep mysql | grep started > /dev/null 2>&1; then
    echo "✅ MySQL服务正在运行"
elif sudo launchctl list | grep mysql > /dev/null 2>&1; then
    echo "✅ MySQL服务正在运行"
else
    echo "❌ MySQL服务未运行，尝试启动..."
    if command -v brew >/dev/null 2>&1; then
        brew services start mysql
        echo "⏳ 等待MySQL启动..."
        sleep 3
    else
        echo "❌ 请手动启动MySQL服务"
        exit 1
    fi
fi

# 检查Go环境
echo ""
echo "🔧 检查Go环境..."
if command -v go >/dev/null 2>&1; then
    echo "✅ Go环境正常: $(go version)"
else
    echo "❌ Go环境未安装"
    exit 1
fi

# 检查依赖
echo ""
echo "📦 检查项目依赖..."
if [ ! -f "go.mod" ]; then
    echo "❌ 未找到go.mod文件"
    exit 1
fi

echo "✅ 项目依赖检查完成"

# 启动应用
echo ""
echo "🚀 启动Go数据库演示应用..."
echo "==============================="
echo "🌐 Web界面: http://localhost:8080"
echo "📖 API接口: http://localhost:8080/api"
echo ""
echo "按 Ctrl+C 停止服务器"
echo ""

go run src/main.go
