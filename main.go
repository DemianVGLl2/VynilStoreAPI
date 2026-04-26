package main

import (
	"github.com/DemianVGLl2/VynilStoreAPI/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Public
	r.POST("/login", handlers.Login)

	// Protected — RequireAuth runs before every handler in this group
	protected := r.Group("/")
	protected.Use(handlers.RequireAuth)
	{
		protected.POST("/logout", handlers.Logout)
		protected.GET("/albums", handlers.GetAlbums)
		protected.GET("/albums/:id", handlers.GetAlbumByID)
		protected.POST("/createAlbum", handlers.CreateAlbum)
		protected.GET("/status", handlers.Status)
	}

	r.Run(":8080")
}
