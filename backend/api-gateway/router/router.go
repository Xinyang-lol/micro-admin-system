package router

import (
	"time"

	"micro-admin-system/backend/api-gateway/client"
	"micro-admin-system/backend/api-gateway/handler"
	"micro-admin-system/backend/api-gateway/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Setup(h *handler.Handler, secret string, redisClient *redis.Client, clients *client.Clients, timeout time.Duration) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	api.POST("/auth/login", h.Login)

	authed := api.Group("")
	authed.Use(middleware.Auth(secret, redisClient))
	authed.POST("/auth/logout", h.Logout)
	authed.GET("/auth/profile", h.Profile)

	users := authed.Group("/users")
	users.GET("", middleware.Permission("sys:user:list", clients, timeout), h.ListUsers)
	users.POST("", middleware.Permission("sys:user:create", clients, timeout), h.CreateUser)
	users.PUT("/:id", middleware.Permission("sys:user:update", clients, timeout), h.UpdateUser)
	users.DELETE("/:id", middleware.Permission("sys:user:delete", clients, timeout), h.DeleteUser)
	users.PUT("/:id/status", middleware.Permission("sys:user:update", clients, timeout), h.UpdateUserStatus)
	users.PUT("/:id/password", middleware.Permission("sys:user:password", clients, timeout), h.ResetUserPassword)
	users.PUT("/:id/roles", middleware.Permission("sys:user:roles", clients, timeout), h.AssignUserRoles)

	roles := authed.Group("/roles")
	roles.GET("", middleware.Permission("sys:role:list", clients, timeout), h.ListRoles)
	roles.POST("", middleware.Permission("sys:role:create", clients, timeout), h.CreateRole)
	roles.PUT("/:id", middleware.Permission("sys:role:update", clients, timeout), h.UpdateRole)
	roles.DELETE("/:id", middleware.Permission("sys:role:delete", clients, timeout), h.DeleteRole)
	roles.PUT("/:id/menus", middleware.Permission("sys:role:menus", clients, timeout), h.AssignRoleMenus)

	menus := authed.Group("/menus")
	menus.GET("/tree", middleware.Permission("sys:menu:list", clients, timeout), h.MenuTree)
	menus.POST("", middleware.Permission("sys:menu:create", clients, timeout), h.CreateMenu)
	menus.PUT("/:id", middleware.Permission("sys:menu:update", clients, timeout), h.UpdateMenu)
	menus.DELETE("/:id", middleware.Permission("sys:menu:delete", clients, timeout), h.DeleteMenu)

	depts := authed.Group("/depts")
	depts.GET("/tree", middleware.Permission("sys:dept:list", clients, timeout), h.DeptTree)
	depts.POST("", middleware.Permission("sys:dept:create", clients, timeout), h.CreateDept)
	depts.PUT("/:id", middleware.Permission("sys:dept:update", clients, timeout), h.UpdateDept)
	depts.DELETE("/:id", middleware.Permission("sys:dept:delete", clients, timeout), h.DeleteDept)

	devices := authed.Group("/devices")
	devices.GET("", middleware.Permission("device:list", clients, timeout), h.ListDevices)
	devices.POST("", middleware.Permission("device:create", clients, timeout), h.CreateDevice)
	devices.PUT("/:id", middleware.Permission("device:update", clients, timeout), h.UpdateDevice)
	devices.DELETE("/:id", middleware.Permission("device:delete", clients, timeout), h.DeleteDevice)
	devices.GET("/statistics", middleware.Permission("device:statistics", clients, timeout), h.DeviceStatistics)

	deviceTypes := authed.Group("/device-types")
	deviceTypes.GET("", middleware.Permission("device:type:list", clients, timeout), h.ListDeviceTypes)
	deviceTypes.POST("", middleware.Permission("device:type:create", clients, timeout), h.CreateDeviceType)
	deviceTypes.PUT("/:id", middleware.Permission("device:type:update", clients, timeout), h.UpdateDeviceType)
	deviceTypes.DELETE("/:id", middleware.Permission("device:type:delete", clients, timeout), h.DeleteDeviceType)

	files := authed.Group("/files")
	files.POST("/upload", middleware.Permission("file:upload", clients, timeout), h.UploadFile)
	files.GET("", middleware.Permission("file:list", clients, timeout), h.ListFiles)
	files.GET("/:id/download", middleware.Permission("file:download", clients, timeout), h.DownloadFile)
	files.DELETE("/:id", middleware.Permission("file:delete", clients, timeout), h.DeleteFile)

	return r
}
