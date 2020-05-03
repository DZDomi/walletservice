package main

import (
	"github.com/DZDomi/walletservice/handlers"
	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine) {

	root := r.Group("/v1")
	{
		wallets := root.Group("/wallets")
		{
			wallets.GET("/:id", handlers.GetWallet)
		}
	}
}
