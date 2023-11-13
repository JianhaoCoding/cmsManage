package main

import (
	"cms/cache"
	conf "cms/config"
	"cms/models"
	"cms/router"
	"fmt"
	"github.com/fvbock/endless"
	"log"
	"strconv"
	"syscall"
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

	// 服务启动端口
	startPortStr := ":" + strconv.Itoa(conf.InitConf().Server.Port)

	// 使用默认router启动
	// r.Run(startPortStr)

	// 使用 endless 来启动服务，而不是用 router.Run()
	server := endless.NewServer(startPortStr, r)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server stopped: %v", err)
	}

	// 更新了代码并且想要应用这些更改时，需要重新编译你的服务
	// go build -o myservice // 调试的时候我们直接在goland下运行项目就好不用编译然后kill -SIGUSR1 <pid>，如果不知道pid直接终端lsof -i:8082
	// kill -SIGUSR1 <pid>
}
