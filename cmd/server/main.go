package main

import (
	"fmt"

	"github.com/bootcamp-go/custom-mocking/cmd/server/handler"
	"github.com/bootcamp-go/custom-mocking/internal/products"
	"github.com/bootcamp-go/custom-mocking/pkg/store"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("%serror al cargar archivo .env %s\n", "\033[31m", "\033[0m")
	}
	db := store.New(store.FileType, "products.json")
	repo := products.NewRepository(db)
	service := products.NewService(repo)
	p := handler.NewProduct(service)

	r := gin.Default()
	pr := r.Group("/products")
	pr.POST("/", p.Store())
	pr.GET("/", p.GetAll())
	r.Run()
}

//go run cmd/server/main.go
