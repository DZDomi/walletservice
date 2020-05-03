package main

import (
	"fmt"
	"github.com/DZDomi/walletservice/clients"
	"github.com/DZDomi/walletservice/models"
	"github.com/gin-gonic/gin"
)

func main() {

	models.InitDB("walletservice:password@(localhost:3307)/walletservice?charset=utf8&parseTime=True&loc=Local", false)
	defer models.DB.Close()

	clients.InitRedis("wallet")
	clients.InitLock()

	clients.InitKafka()

	gin.SetMode("debug")
	r := gin.Default()

	routes(r)

	_ = r.Run(fmt.Sprintf(":%d", 8081))
}
