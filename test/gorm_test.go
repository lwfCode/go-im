package test

import (
	"fmt"
	"ginchat/models"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGorm(t *testing.T) {
	dns := "root:skyinfor2016@tcp(192.168.6.185:38)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic("connect mysql error!")
	}

	// db.AutoMigrate(&models.UserBasic{})

	message := &models.Contact{}

	message.UserId = 1
	message.TargetId = 2
	message.Type = 1
	message.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	message.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	db.Create(message)

	db.First(message)
	fmt.Println(message, "======first")

	// db.Model(message).Update("Content", "123456")
}
