package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func main() {
	// 连接数据库
	dsn := "root:12345678@tcp(127.0.0.1:3306)/poker?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 解决查表的时候会自动添加复数的问题, 例如userInfo 变成了userInfos
			SingularTable: true,
		},
	})
	fmt.Println(db)
	fmt.Println(err)

	sqlDB, err := db.DB()
	if err != nil {
		// 处理错误，例如打印错误信息或者返回错误给调用者
		fmt.Println("Failed to get DB instance:", err)
		return
	}
	// setMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// setMaxIdleConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)

	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 10s

	type UserInfo struct {
		gorm.Model
		UserName  string `gorm:"type:varchar(20);not null;" json:"userName" binding:"required"`
		AvatarUrl string `gorm:"type:varchar(20);not null;" json:"avatarUrl" binding:"required"`
		OpenId    string `gorm:"type:varchar(20);not null" json:"openId" binding:"required"`
	}
	db.AutoMigrate(&UserInfo{})

	// 接口
	r := gin.Default()
	// 添加userInfo
	r.POST("setUserAvatar", func(c *gin.Context) {
		var userInfo UserInfo
		if err := c.ShouldBindJSON(&userInfo); err != nil {
			c.JSON(200, gin.H{
				"message": "error",
				"data":    gin.H{},
			})
		}
		db.Create(&userInfo)
		c.JSON(200, gin.H{
			"message": "success",
			"data":    userInfo,
		})
	})
	r.GET("/getUserInfo", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "poker",
		})
	})
	// 端口号
	PORT := "8000"
	r.Run(":" + PORT)
}
