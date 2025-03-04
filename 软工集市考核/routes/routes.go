package routes

import (
	"vote-demo/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, User-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 用户相关路由
	userRoutes := r.Group("/api/users")
	{
		userRoutes.POST("", controllers.CreateUser)
		userRoutes.GET("", controllers.ListUsers)
		userRoutes.GET("/:id", controllers.GetUser)
		userRoutes.GET("/username/:username", controllers.GetUserByUsername)
		userRoutes.GET("/:id/stats", controllers.GetUserStats)
	}

	// 投票相关路由
	pollRoutes := r.Group("/api/polls")
	{
		pollRoutes.POST("", controllers.CreatePoll)
		pollRoutes.GET("", controllers.ListPolls)
		pollRoutes.GET("/:id", controllers.GetPoll)
		pollRoutes.PUT("/:id", controllers.UpdatePoll)
		pollRoutes.DELETE("/:id", controllers.DeletePoll)
		pollRoutes.GET("/:id/results", controllers.GetPollResults)
		pollRoutes.GET("/:id/stats", controllers.GetPollStats)

		// 选项相关路由
		pollRoutes.POST("/:id/options", controllers.AddOption)
		pollRoutes.PUT("/:id/options/:option_id", controllers.UpdateOption)
		pollRoutes.DELETE("/:id/options/:option_id", controllers.DeleteOption)

		// 投票操作路由
		pollRoutes.POST("/:id/vote", controllers.CastVote)
		pollRoutes.GET("/:id/user-votes", controllers.GetUserVotes)

		// 评论相关路由
		pollRoutes.POST("/:id/comments", controllers.AddComment)
		pollRoutes.GET("/:id/comments", controllers.GetPollComments)
		pollRoutes.PUT("/:id/comments/:comment_id", controllers.UpdateComment)
		pollRoutes.DELETE("/:id/comments/:comment_id", controllers.DeleteComment)
	}

	// 统计和分析路由
	statsRoutes := r.Group("/api/stats")
	{
		statsRoutes.GET("/trending", controllers.GetTrendingPolls)
	}

	return r
} 