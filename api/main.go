package main

import "github.com/gin-gonic/gin"

func main() {
	initDB()

	r := gin.Default()

	r.POST("/jobs", createJob)
	r.GET("/jobs/:id", getJob)

	r.Run(":8080")
}
