package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yurin-kami/PackChann/database"
	"github.com/yurin-kami/PackChann/middlewares"
	"github.com/yurin-kami/PackChann/models"
	"github.com/yurin-kami/PackChann/routes"
	"github.com/yurin-kami/PackChann/utils"
)

func main() {
	fmt.Println("PackChann System On~")

	// 0. 初始化 Snowflake
	if err := utils.InitSnowflake(1); err != nil {
		log.Fatalf("无法初始化 Snowflake: %v", err)
	}

	// 1. 加载配置
	cfg, err := models.LoadConfig()
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 2. 连接数据库
	db, err := database.NewConnection(context.Background(), cfg.Database)
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}

	router := gin.Default()

	// 添加 CORS 中间件
	router.Use(middlewares.CORSMiddleware())

	// 3. 注册路由
	// 未受保护路由 (登录/注册)
	routes.UnprotectedRoutes(db, router, cfg.JWT)

	// 受保护路由 (业务逻辑)
	routes.ProtectedRoutes(db, router, cfg.JWT)

	router.Run(":8088")
}
