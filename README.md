# Go Web API 学习项目

一个基于 Go + Gin + GORM + MySQL 的现代化 Web API 学习项目，包含用户管理、留言板功能和完整的前端界面。

![Go Version](https://img.shields.io/badge/Go-1.24.5-blue)
![Gin](https://img.shields.io/badge/Gin-1.10.1-green)
![GORM](https://img.shields.io/badge/GORM-1.30.0-orange)
![MySQL](https://img.shields.io/badge/MySQL-8.0-blue)

## 🎯 项目概述

这是一个功能完整的 Go Web 开发学习项目，实现了：

- ✅ RESTful API 设计
- ✅ 动态项目路径管理
- ✅ 智能数据库连接和迁移
- ✅ 用户 CRUD 操作
- ✅ 留言板功能（支持IP地址获取）
- ✅ 现代化 Web 前端界面
- ✅ 智能服务器重启脚本
- ✅ 跨域 CORS 支持
- ✅ 健康检查和监控

## 🎆 效果展示
<img src="https://raw.githubusercontent.com/OPTJoker/goLearn/refs/heads/main/webScreenShot.png" alt="效果演示" width="375">

## 🚀 快速开始

### 环境要求

- Go 1.24.5+
- MySQL 8.0+
- curl（用于API测试）

### 镜像配置
<img width="403" height="154" alt="image" src="https://github.com/user-attachments/assets/b3a7ae3d-4afa-486d-b0c9-4cd05191c6b3" />

**神奇的大陆用户，拉依赖库的时候可能会遇到类似上面的报错，所以你需要配置如下镜像：**
```bash
# 设置七牛云镜像源（永久生效）
go env -w GOPROXY=https://goproxy.cn,direct

# 或设置阿里云镜像源
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 验证配置
go env | grep GOPROXY
```

### 安装依赖

```bash
# 克隆项目
git clone https://github.com/OPTJoker/goLearn.git
cd goLearn

# 安装Go依赖
go mod tidy
```

### 配置数据库

MySQL配置（默认配置）：
- 主机: localhost
- 端口: 3306
- 用户: root
- 密码: （空）
- 数据库: xldb_webapi

### 启动服务

```bash
# 方法1: 直接启动
go run src/main.go

# 方法2: 使用智能重启脚本（推荐）
chmod +x script/restart_server.sh
./script/restart_server.sh
```

服务启动后，访问：
- 🌐 Web界面: http://localhost:8080
- 📡 API接口: http://localhost:8080/api/
- 📊 数据库状态: http://localhost:8080/api/database/status

## 📁 项目结构

```
xlgo/
├── src/
│   └── main.go              # 主程序入口（支持动态路径）
├── web/
│   ├── index.html           # 现代化Web界面
│   └── script.js            # 前端JavaScript逻辑
├── util/
│   ├── util.go              # IP获取等工具函数
│   └── project.go           # 项目路径管理工具
├── script/
│   └── restart_server.sh    # 智能重启脚本
├── go.mod                   # Go模块配置
├── go.sum                   # 依赖锁定文件
└── README.md                # 项目说明
```

## 🛠️ 核心技术栈

| 技术 | 版本 | 用途 | 特点 |
|------|------|------|------|
| Go | 1.24.5 | 后端语言 | 高性能、并发支持 |
| Gin | 1.10.1 | Web框架 | 轻量级、高性能路由 |
| GORM | 1.30.0 | ORM框架 | 自动迁移、类型安全 |
| MySQL | 8.0+ | 数据库 | ACID事务、高可用 |
| HTML5/JS | - | 前端界面 | 现代化用户体验 |

## 📡 API 接口文档

### 数据库管理 API

| 方法 | 路径 | 描述 | 状态 |
|------|------|------|------|
| POST | `/api/database/create` | 创建数据库 | ✅ |
| POST | `/api/database/connect` | 连接数据库（自动迁移） | ✅ |
| GET | `/api/database/status` | 数据库状态检查 | ✅ |

#### 🔗 创建数据库
```bash
curl -X POST http://localhost:8080/api/database/create \
  -H "Content-Type: application/json" \
  -d '{
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "",
    "dbname": "xldb_webapi"
  }'
```

#### 🔌 连接数据库
```bash
curl -X POST http://localhost:8080/api/database/connect \
  -H "Content-Type: application/json" \
  -d '{
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "",
    "dbname": "xldb_webapi"
  }'
```

### 用户管理 API

| 方法 | 路径 | 描述 | 功能特点 |
|------|------|------|----------|
| POST | `/api/users` | 创建用户 | 数据验证、唯一性检查 |
| GET | `/api/users` | 获取所有用户 | 分页支持 |
| GET | `/api/users/:id` | 获取指定用户 | 错误处理 |
| PUT | `/api/users/:id` | 更新用户 | 部分更新支持 |
| DELETE | `/api/users/:id` | 删除用户 | 软删除 |

#### 👤 创建用户
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "age": 25
  }'
```

### 留言板 API

| 方法 | 路径 | 描述 | 特色功能 |
|------|------|------|----------|
| POST | `/api/addContent` | 添加留言 | 自动IP地址获取 |
| GET | `/api/getAllContent` | 获取所有留言 | 时间排序 |
| DELETE | `/api/removeContent/:msg_id` | 删除留言 | 权限控制 |

#### 💬 添加留言
```bash
curl -X POST http://localhost:8080/api/addContent \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "用户123",
    "content": "这是一条测试留言"
  }'
```

**注意**: `user_ip` 字段会自动填充，支持：
- 本地开发显示："本地访问(localhost)"
- 生产环境显示：真实IP地址
- 代理环境：从HTTP头提取真实IP

## 🗄️ 数据库结构

### users 表
```sql
CREATE TABLE `users` (
  `id` bigint unsigned AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `email` varchar(100) UNIQUE,
  `age` int,
  `created_at` datetime(3),
  `updated_at` datetime(3),
  PRIMARY KEY (`id`)
);
```

### msg_contents 表
```sql
CREATE TABLE `msg_contents` (
  `msg_id` bigint unsigned AUTO_INCREMENT,
  `user_id` text NOT NULL,
  `user_ip` text NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime(3),
  PRIMARY KEY (`msg_id`)
);
```

## 🔧 开发工具与脚本

### 智能重启服务器
```bash
# 推荐方式：使用智能重启脚本
./script/restart_server.sh
```

智能重启脚本特点：
- ✅ 优雅停止现有进程
- ✅ 健康检查确认启动成功
- ✅ 详细的状态反馈
- ✅ 错误处理和回滚

### 项目路径管理

项目支持多种路径获取方式：

```bash
# 使用环境变量（生产环境推荐）
export PROJECT_ROOT=/path/to/your/project
go run src/main.go

# 自动检测（开发环境）
cd /path/to/project
go run src/main.go

# 从src目录运行
cd src
go run main.go
```

### 手动操作
```bash
# 查看服务状态
lsof -i :8080

# 手动停止服务
pkill -f "go run"

# 查看日志（如果启用）
tail -f server.log
```

## 🌐 Web界面功能

现代化的Web界面包含：

- 📊 **数据库管理面板**: 创建、连接、状态检查
- 👥 **用户管理界面**: CRUD操作、实时反馈
- 💬 **留言板系统**: 发布留言、查看历史
- 🎨 **响应式设计**: 支持桌面和移动设备
- ⚡ **实时更新**: Ajax异步操作

## 📈 功能特性

### ✅ 已实现功能
- 🏗️ RESTful API设计
- 🔄 自动数据库迁移
- 🌍 CORS跨域支持
- 📝 统一错误处理和日志
- 🎯 智能IP地址获取
- 📱 响应式Web界面
- 🔍 健康检查接口
- 🚀 智能重启脚本
- 📂 动态项目路径管理

### 🔄 计划功能(并没有真的在计划)
- 🔐 JWT用户认证和授权
- ✅ 数据验证和输入校验
- ⚡ API限流和Redis缓存
- 🧪 单元测试和集成测试
- 🐳 Docker容器化部署
- 🔄 CI/CD自动化流水线
- 📊 性能监控和指标收集

## 🎓 学习要点

这个项目涵盖了以下Go Web开发的核心概念：

### 后端开发
1. **Web框架使用** - Gin框架的路由、中间件、请求处理
2. **数据库操作** - GORM的模型定义、自动迁移、CRUD操作
3. **RESTful设计** - HTTP方法、状态码、资源设计原则
4. **错误处理** - 统一的错误响应格式和日志记录

### 系统设计
5. **配置管理** - 环境变量、数据库配置、项目路径
6. **项目结构** - 模块化设计、代码组织、依赖管理
7. **部署运维** - 脚本自动化、健康检查、进程管理

### 前端交互
8. **API设计** - 接口规范、数据格式、错误处理
9. **前后端分离** - Ajax请求、数据绑定、用户体验

## 🛡️ 故障排除

### 常见问题及解决方案

#### 1. 端口被占用
```bash
# 查看占用进程
lsof -i :8080

# 使用智能重启脚本（推荐）
./script/restart_server.sh

# 手动终止进程
kill -9 <PID>
```

#### 2. 数据库连接失败
- ✅ 检查MySQL服务是否启动
- ✅ 确认连接参数（host、port、user、password）
- ✅ 检查用户权限和数据库是否存在
- ✅ 使用API先创建数据库

#### 3. 依赖问题
```bash
# 重新下载依赖
go mod tidy

# 清理模块缓存
go clean -modcache

# 验证模块
go mod verify
```

#### 4. 表结构问题
- ✅ 系统会自动执行数据库迁移
- ✅ 通过API连接数据库时自动创建缺失字段
- ✅ 检查日志输出确认迁移状态

#### 5. 静态文件404
- ✅ 项目自动检测根目录
- ✅ 支持环境变量 `PROJECT_ROOT` 配置
- ✅ 查看启动日志确认路径

### 调试模式

```bash
# 开启详细日志
export GIN_MODE=debug
go run src/main.go

# 检查项目路径
export PROJECT_ROOT=/your/project/path
go run src/main.go
```

## 🚀 部署指南

### 开发环境
1. 克隆代码并安装依赖
2. 启动MySQL服务
3. 运行 `go run src/main.go`
4. 访问 http://localhost:8080

### 生产环境
1. 设置环境变量 `PROJECT_ROOT`
2. 配置MySQL连接参数
3. 使用智能重启脚本
4. 配置反向代理（Nginx）

## 🤝 贡献指南

欢迎贡献代码和提出建议！

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request


## 📞 联系方式

- 作者: Sharon
- 项目链接: [GitHub Repository](https://github.com/OPTJoker/goLearn)

## 📝 许可证
本项目不采取任何许可证，随便造

## 🙏 致谢

感谢以下开源项目：
- [Gin](https://github.com/gin-gonic/gin) - 高性能HTTP Web框架
- [GORM](https://github.com/go-gorm/gorm) - 强大的ORM库
- [MySQL Driver](https://github.com/go-gorm/mysql) - MySQL数据库驱动

---

