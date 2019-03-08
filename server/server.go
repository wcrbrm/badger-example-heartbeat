package main

import (
	_ "github.com/dgraph-io/badger"
	gin     "github.com/gin-gonic/gin"
	ginprom "github.com/Depado/ginprom"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"fmt"
)

func main() {
	fmt.Println("Starting")

	r := gin.Default()
	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("gin"), 
		ginprom.Path("/metrics"), 
	)
	r.Use(p.Instrument())

	r.POST("/heartbeat", func(c *gin.Context) {})
	r.GET("/active", func(c *gin.Context) {})
        r.GET("/all",   func(c *gin.Context) {})
        r.POST("/search",  func(c *gin.Context) {})
	r.Run("0.0.0.0:8092")
}
