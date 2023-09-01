package proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Start() {
	// 修改gin模式
	gin.SetMode(gin.ReleaseMode)
	// 创建路由
	router := gin.Default()

	router.POST("/payAuthRes", func(c *gin.Context) {
		log.Printf("可以了")
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	fmt.Println("proxy start")

	router.Run(":8866")
}
