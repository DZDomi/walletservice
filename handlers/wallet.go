package handlers

import (
	"github.com/DZDomi/walletservice/clients"
	"github.com/DZDomi/walletservice/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetWallet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lock, err := clients.GetLock(string(id), time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer clients.ReleaseLock(lock)

	wallet := new(models.Wallet)
	models.DB.Where("id = ?", id).First(&wallet)

	if wallet.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}
