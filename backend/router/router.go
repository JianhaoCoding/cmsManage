package router

import (
	"cms/controllers"
	"cms/controllers/adminer"
	"cms/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置 CORS
	// r.Use(cors.Default())	// 默认配置
	// gin 的 cors 中间件确实需要 AllowOrigins 字段。如果您不提供 AllowOrigins，则需要至少设置 AllowAllOrigins 为 true，或者提供 AllowOriginFunc 函数。
	// 希望从任意源接受请求，并且这些请求不需要携带认证信息，您可以将 AllowAllOrigins 设置为 true 并且将 AllowCredentials 设置为 false
	r.Use(cors.New(cors.Config{
		// 允许的源
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}, // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},                  // 允许的头部字段
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // 当你设置 AllowCredentials: true 时，你不能将 Access-Control-Allow-Origin 设置为 "*"（允许任何域），因为这被认为是不安全的。你必须指定一个确切的、可信的域名作为请求的来源
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "http://localhost:8080" || origin == "http://example.com"
		//},
		MaxAge: 12 * time.Hour,
	}))

	// 允许指定的源
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://127.0.0.1:8080"}
	//config.AllowMethods = []string{"GET", "POST", "PUT"}
	//config.AllowHeaders = []string{"Origin"}
	//config.ExposeHeaders = []string{"Content-Length"}
	//config.AllowCredentials = true
	//r.Use(cors.New(config))

	// JWT保护路由的3种方式
	// 1 r.POST("/menurules", middleware.TokenAuth(), adminer.MenuAuths)
	// 2 使用路由分组保护如下：
	// auth := r.Group("/auth")
	// auth.Use(middleware.TokenAuth()) {
	//		auth.GET("/menurules", adminer.MenuAuths)
	//}
	// 3 直接在JWT中间件中过滤不需要过滤的接口，设置JWT为全局中间件
	//r.Use(middleware.TokenAuth())
	//func TokenAuth() gin.HandlerFunc {
	//	return func(c *gin.Context) {
	//	// 跳过某些路由
	//	if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/public" {
	//	c.Next()
	//	return
	//}
	//	// JWT验证逻辑...

	// 选择第三种，好控制方便重写什么的
	r.Use(middleware.TokenAuth())
	// 登录操作
	r.POST("login", controllers.Login)

	// 管理员操作
	r.GET("/menurules", adminer.MenuAuths)                         // 用户Menu列表接口
	r.GET("/adminer/list", adminer.AdminerList)                    // 管理员列表
	r.GET("/adminer/show/:adminer_id", adminer.AdminerShow)        // 管理员信息
	r.POST("/adminer/addadminer", adminer.AddAdminer)              // 添加管理员
	r.POST("/adminer/switchstatus", adminer.UpdateAdminerStatus)   // 管理员状态开关
	r.DELETE("/adminer/delete/:adminer_id", adminer.DeleteAdminer) // 删除管理员
	r.POST("/adminer/add", adminer.AddAdminer)                     // 添加管理员
	r.POST("/adminer/edit", adminer.UpdateAdminer)                 // 编辑管理员
	r.GET("/adminer/resetpass", adminer.ResetPass)                 // 重置管理员密码

	// 管理组操作
	r.GET("/adminer/group/list", adminer.AdminerGroupList)           // 管理员组列表
	r.GET("/adminer/group/show/:group_id", adminer.AdminerGroupShow) // 管理员组信息
	r.POST("/adminer/group/save", adminer.SaveAdminerGroup)          // 保存管理员信息（添加+更新）

	// 管理权限操作
	r.GET("/adminer/auth/list", adminer.AuthList)                  // 权限管理列表
	r.GET("/adminer/auth/show/:auth_id", adminer.AuthShow)         // 权限详情
	r.POST("/adminer/auth/switchstatus", adminer.SwitchAuthStatus) // 权限状态开关
	r.POST("/adminer/auth/switchmenu", adminer.SwitchMenuShow)     // 权限状态开关
	r.POST("/adminer/auth/save", adminer.SaveAuth)                 // 保存权限
	r.DELETE("/adminer/auth/delete/:auth_id", adminer.DeleteAuth)  // 删除权限

	return r
}
