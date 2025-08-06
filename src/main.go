package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"xlgo/util"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 数据库连接实例
var db *gorm.DB

// 获取项目根目录
func getProjectRoot() string {
	// 方法1: 检查环境变量 PROJECT_ROOT
	if projectRoot := os.Getenv("PROJECT_ROOT"); projectRoot != "" {
		fmt.Printf("📌 使用环境变量 PROJECT_ROOT: %s\n", projectRoot)
		return projectRoot
	}

	// 方法2: 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("无法获取当前工作目录: %v", err)
		return "."
	}

	// 如果从src目录运行，则向上一级到项目根目录
	if filepath.Base(wd) == "src" {
		return filepath.Dir(wd)
	}

	// 如果已经在项目根目录，直接返回
	return wd
}

// 获取Web目录路径
func getWebDir() string {
	projectRoot := getProjectRoot()
	return filepath.Join(projectRoot, "web")
}

// 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Email     string    `json:"email" gorm:"size:100;uniqueIndex"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MsgContent struct {
	MsgID     uint      `json:"msg_id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"type:text;not null"`
	UserIP    string    `json:"user_ip" gorm:"type:text;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}

// 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// API响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var allDBTables = []any{
	&User{},
	&MsgContent{},
}

func validateDatabase() error {
	// 检查所有必要的表是否存在
	for _, model := range allDBTables {
		if !db.Migrator().HasTable(model) {
			return fmt.Errorf("表 %T 不存在", model)
		}
	}
	return nil
}

// 初始化数据库连接
func initDatabase(config DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)
	print(">>> 连接指定数据库参数: ", dsn)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 强制自动迁移表结构以确保所有字段存在
	fmt.Println("开始自动迁移表结构...")
	if err := db.AutoMigrate(allDBTables...); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	} else {
		fmt.Println("表结构迁移成功！")
	}

	fmt.Println("数据库连接成功！")
	return nil
}

// 创建数据库
func createDatabase(config DatabaseConfig) error {
	// 连接到MySQL服务器（不指定数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port)
	print(">>> 连接数据库参数: ", dsn)

	tempDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %v", err)
	}

	// 创建数据库
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", config.DBName)
	if err := tempDB.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	fmt.Printf("数据库 %s 创建成功！\n", config.DBName)
	return nil
}

// HTTP处理函数

// 创建数据库接口
func createDatabaseHandler(c *gin.Context) {
	var config DatabaseConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if err := createDatabase(config); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: fmt.Sprintf("数据库 %s 创建成功", config.DBName),
	})
}

// 连接数据库接口
func connectDatabaseHandler(c *gin.Context) {
	var config DatabaseConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if err := initDatabase(config); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "数据库连接成功",
	})
}

// 创建用户接口
func createUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if db == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "数据库未连接",
		})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "创建用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "用户创建成功",
		Data:    user,
	})
}

func addContent(c *gin.Context) {
	var msgContent MsgContent
	if err := c.ShouldBindJSON((&msgContent)); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 填充ip地址
	msgContent.UserIP = util.GetClientIP(c)

	if db == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "数据库未连接",
		})
		return
	}
	if err := db.Create(&msgContent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "评论失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "评论成功",
		Data:    msgContent,
	})

}

func getAllContent(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "数据库未连接",
		})
		return
	}

	var msgContents []MsgContent
	if err := db.Find(&msgContents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "查询评论失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "查询成功",
		Data:    msgContents,
	})
}

// 获取所有用户接口
func getUsersHandler(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "数据库未连接",
		})
		return
	}

	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "查询用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "查询成功",
		Data:    users,
	})
}

// 根据ID获取用户接口
func getUserByIDHandler(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "数据库未连接",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID参数错误",
		})
		return
	}

	var user User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "查询用户失败: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "查询成功",
		Data:    user,
	})
}

// 更新用户接口
func updateUserHandler(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "数据库未连接",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID参数错误",
		})
		return
	}

	var user User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Message: "用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "查询用户失败: " + err.Error(),
			})
		}
		return
	}

	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	if err := db.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "更新用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "用户更新成功",
		Data:    user,
	})
}

// 删除用户接口
func deleteUserHandler(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "数据库未连接",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "ID参数错误",
		})
		return
	}

	if err := db.Delete(&User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "删除用户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "用户删除成功",
	})
}

// 删除留言接口
func deleteContentHandler(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "数据库未连接",
		})
		return
	}
	msgId, err := strconv.Atoi(c.Param("msg_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "msg_id参数错误",
		})
		return
	}
	if err := db.Delete(&MsgContent{}, msgId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "删除留言失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "留言删除成功",
	})
}

// 数据库状态检查接口
func databaseStatusHandler(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Message: "数据库状态",
			Data: map[string]interface{}{
				"connected": false,
				"status":    "未连接",
			},
		})
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "获取数据库连接失败: " + err.Error(),
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Message: "数据库状态",
			Data: map[string]interface{}{
				"connected": false,
				"status":    "连接失败",
				"error":     err.Error(),
			},
		})
		return
	}

	stats := sqlDB.Stats()
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "数据库状态",
		Data: map[string]interface{}{
			"connected":        true,
			"status":           "已连接",
			"open_connections": stats.OpenConnections,
			"in_use":           stats.InUse,
			"idle":             stats.Idle,
		},
	})
}

func main() {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 数据库管理接口
	api := r.Group("/api")
	{
		// 数据库操作
		api.POST("/database/create", createDatabaseHandler)
		api.POST("/database/connect", connectDatabaseHandler)
		api.GET("/database/status", databaseStatusHandler)

		// 用户CRUD操作
		api.POST("/users", createUserHandler)
		api.GET("/users", getUsersHandler)
		api.GET("/users/:id", getUserByIDHandler)
		api.PUT("/users/:id", updateUserHandler)
		api.DELETE("/users/:id", deleteUserHandler)

		// 留言板操作
		api.POST("/addContent", addContent)
		api.GET("/getAllContent", getAllContent)
		api.DELETE("/removeContent/:msg_id", deleteContentHandler)
	}

	// 静态文件服务 - 使用动态路径
	webDir := getWebDir()
	r.Static("/static", webDir)

	fmt.Printf("📁 项目根目录: %s\n", getProjectRoot())
	fmt.Printf("🌐 Web目录: %s\n", webDir)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	fmt.Println("🚀 服务器启动在 http://localhost:8080")
	fmt.Println("📱 Web界面: http://localhost:8080")
	fmt.Println("📖 API文档:")
	fmt.Println("  POST /api/database/create  - 创建数据库")
	fmt.Println("  POST /api/database/connect - 连接数据库")
	fmt.Println("  GET  /api/database/status  - 数据库状态")
	fmt.Println("  POST /api/users           - 创建用户")
	fmt.Println("  GET  /api/users           - 获取所有用户")
	fmt.Println("  GET  /api/users/:id       - 获取指定用户")
	fmt.Println("  PUT  /api/users/:id       - 更新用户")
	fmt.Println("  DELETE /api/users/:id     - 删除用户")
	fmt.Println("  DELETE /api/users/:id     - 删除用户")

	fmt.Println("  POST /api/addContent   - 添加留言")
	fmt.Println("  GET  /api/getAllContent   - 获取所有留言")

	log.Fatal(r.Run(":8080"))
}
