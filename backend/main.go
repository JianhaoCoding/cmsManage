package main

import (
	"cms/cache"
	conf "cms/config"
	"cms/models"
	"cms/router"
	"fmt"
	"log"
	"strconv"
)

func init() {
	// 初始化全局配置
	conf.InitConf()

	// 初始化数据库
	dataBaseConf := conf.InitConf().Database
	dbClentStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dataBaseConf.User,
		dataBaseConf.Password,
		dataBaseConf.Local,
		dataBaseConf.Port,
		dataBaseConf.Dbname,
	)
	err := models.InitDB(dbClentStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}

	// 初始化redis
	redisClientErr := cache.InitRedis()
	if redisClientErr != nil {
		log.Fatalf("unable to connect to Redis: %v", redisClientErr.Error())
		return
	}
}

func main() {
	// 路由加载 & 运行环境
	r := router.SetupRouter()
	r.Run(":" + strconv.Itoa(conf.InitConf().Server.Port))
}
