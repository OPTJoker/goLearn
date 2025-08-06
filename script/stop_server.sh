#!/bin/bash
# smart_restart.sh

APP_NAME="go_server"
PORT=8080
HEALTH_CHECK_URL="http://localhost:$PORT/api/database/status"
MAX_WAIT_TIME=10

echo "🔄 停止服务器 $APP_NAME..."

# 函数：检查服务器是否响应
check_server() {
    curl -s -f "$HEALTH_CHECK_URL" > /dev/null
    return $?
}

# 优雅停止现有进程
echo "🛑 正在停止现有服务..."
PID=$(lsof -ti:$PORT)
if [ ! -z "$PID" ]; then
    echo "📍 发现进程 $PID，发送终止信号..."
    kill $PID
    
    # 等待优雅停止
    sleep 3
    
    # 如果还在运行，强制终止
    if kill -0 $PID 2>/dev/null; then
        echo "⚡ 强制终止进程..."
        kill -9 $PID
    fi
    echo "✅ 旧进程已停止"
fi