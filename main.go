package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Product struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Price       float64    `json:"price"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "Products"
}

type ProductCreation struct {
	Id          int     `json:"-" gorm:"column:id;"`
	Title       string  `json:"title" gorm:"column:title;"`
	Description string  `json:"description" gorm:"column:description;"`
	Status      string  `json:"status" gorm:"column:status;"`
	Price       float64 `json:"price" gorm:"column:price;"`
}

func (ProductCreation) TableName() string {
	return Product{}.TableName()
}

type ProductUpdate struct {
	Title       *string  `json:"title" gorm:"column:title;"`
	Description *string  `json:"description" gorm:"column:description;"`
	Status      *string  `json:"status" gorm:"column:status;"`
	Price       *float64 `json:"price" gorm:"column:price;"`
}

func (ProductUpdate) TableName() string {
	return Product{}.TableName()
}

func main() {

	dsn := os.Getenv("DB_CONN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(db)

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		products := v1.Group("/products")
		{
			products.GET("")
			products.GET("/:id", GetProductById(db))
			products.POST("", CreateProduct(db))
			products.PATCH("/:id", UpdateProduct(db))
			products.DELETE("/:id", DeleteProduct(db))
		}
	}

	if err := r.Run(); err != nil {
		panic(err)
	}
}

func CreateProduct(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data ProductCreation

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data.Id,
		})
	}
}

func GetProductById(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, ok := strconv.Atoi(c.Param("id"))

		if ok != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ok.Error(),
			})

			return
		}

		var data Product

		if err := db.Where("id = ?", id).Find(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}
}

func UpdateProduct(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, ok := strconv.Atoi(c.Param("id"))

		if ok != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ok.Error(),
			})

			return
		}

		var data ProductUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}

func DeleteProduct(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, ok := strconv.Atoi(c.Param("id"))

		if ok != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ok.Error(),
			})

			return
		}

		if err := db.Table(Product{}.TableName()).Where("id = ?", id).Delete(&Product{}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": ok.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})
	}
}
