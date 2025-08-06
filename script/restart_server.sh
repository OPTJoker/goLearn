#!/bin/bash
# smart_restart.sh

APP_NAME="go_server"
PORT=8080
HEALTH_CHECK_URL="http://localhost:$PORT/api/database/status"
MAX_WAIT_TIME=10

echo "🔄 智能重启 $APP_NAME..."

# 函数：检查服务器是否响应
check_server() {
    curl -s -f "$HEALTH_CHECK_URL" > /dev/null
    return $?
}

# 函数：等待服务器启动
wait_for_server() {
    local wait_time=0
    while [ $wait_time -lt $MAX_WAIT_TIME ]; do
        if check_server; then
            return 0
        fi
        echo "⏳ 等待服务器启动... ($wait_time/$MAX_WAIT_TIME)"
        sleep 1
        wait_time=$((wait_time + 1))
    done
    return 1
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

# 启动新服务器
echo "🚀 启动新服务器..."
# nohup go run src/main.go > server.log 2>&1 & 
go run src/main.go
NEW_PID=$!

# 等待服务器启动
if wait_for_server; then
    echo "✅ 服务器重启成功！"
    echo "🆔 新进程ID: $NEW_PID"
    echo "🌐 访问地址: http://localhost:$PORT"
    # echo "📝 日志文件: server.log"
else
    echo "❌ 服务器启动失败"
    # echo "📝 查看日志: tail -f server.log"
    exit 1
fi