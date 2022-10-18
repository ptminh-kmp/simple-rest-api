package main

import (
	"log"
	"net/http"
	"os"
	"simple-rest-api/component"
	"simple-rest-api/modules/restaurant/restauranttransport/ginrestaurant"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DBConnectionStr")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	if err := runService(db); err != nil {
		log.Fatalln(err)
	}
}

func runService(db *gorm.DB) error {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	appCtx := component.NewAppContext(db)

	restaurants := r.Group("/restaurants")
	{
		restaurants.POST("", ginrestaurant.CreateRestaurant(appCtx))

		restaurants.GET("/:id", ginrestaurant.GetRestaurant(appCtx))

		restaurants.GET("", ginrestaurant.ListRestaurant(appCtx))

		restaurants.PATCH("/:id", func(ctx *gin.Context) {
			id, err := strconv.Atoi(ctx.Param("id"))

			if err != nil {
				ctx.JSON(401, gin.H{"error": err.Error()})

				return
			}

			var data RestaurantUpdate

			if err := ctx.ShouldBind(&data); err != nil {
				ctx.JSON(401, gin.H{"error": err.Error()})
				return
			}

			if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
				ctx.JSON(401, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{"ok": 1})
		})

		restaurants.DELETE("/:id", func(ctx *gin.Context) {
			id, err := strconv.Atoi(ctx.Param("id"))

			if err != nil {
				ctx.JSON(401, gin.H{"error": err.Error()})

				return
			}

			if err := db.Table(Restaurant{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
				ctx.JSON(401, gin.H{"error": err.Error()})

				return
			}

			ctx.JSON(http.StatusOK, gin.H{"ok": 1})

		})
	}

	return r.Run()
}

type Restaurant struct {
	Id   int    `json:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name;"`
	Addr string `json:"address" gorm:"column:addr;"`
}

func (Restaurant) TableName() string {
	return "restaurants"
}

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name;"`
	Addr *string `json:"address" gorm:"column:addr;"`
}

func (RestaurantUpdate) TableName() string {
	return Restaurant{}.TableName()
}
