package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tokha04/blogging-platform-api/controllers"
)

func Routes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/posts", controllers.CreateBlog())
	incomingRoutes.PATCH("/posts/:id", controllers.UpdateBlog())
	incomingRoutes.DELETE("/posts/:id", controllers.DeleteBlog())
	incomingRoutes.GET("/posts/:id", controllers.GetBlog())
	incomingRoutes.GET("/posts", controllers.GetBlogs())
}
