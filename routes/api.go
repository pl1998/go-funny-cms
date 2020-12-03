package routes

import (
	"github.com/gin-gonic/gin"
	Admin "gocms/app/http/admin/controllers"
	"gocms/app/http/middleware"
)

// 路由注册
func RegisterApiRoutes(router *gin.Engine) {
	// 后台项目
	apiRouter := router.Group("api")

	// 不需要登陆的路由
	authController := new(Admin.AuthController)
	apiRouter.POST("/admin/login", authController.Login)
	apiRouter.POST("/admin/register", authController.Register)

	// 需要鉴权的路由
	apiRouter.Use(middleware.Authenticate)
	{
		apiRouter.GET("/admin/me", authController.Me)
		apiRouter.POST("/admin/logout", authController.Logout)

		AdminController := new(Admin.AdminController)
		apiRouter.GET("/admin/list", AdminController.Index)
		apiRouter.POST("/admin/add", AdminController.Add)

		ToolController := new(Admin.ToolController)
		apiRouter.GET("/admin/qiniu", ToolController.Qiniu)
	}
}
