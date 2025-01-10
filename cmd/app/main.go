package main

import (
	"awesomeProject1/config"
	"awesomeProject1/internal/db"
	"awesomeProject1/internal/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

// Configurations
//const (
//	reportEndpoint = "/generate-report"
//	loginEndpoint  = "/login"
//)

// Inventory represents the database table
//type Inventory struct {
//	ItemID   int `gorm:"primaryKey"`
//	Quantity int
//	Price    float64
//}

func main() {

	config.InitLogging()
	db.InitDB()
	defer db.Client.Close()
	router := gin.Default()
	handlers.SetupRoutes(router)
	log.Println("Starting server on 8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

	// Insert a new inventory item
	//newItem, err := client.Inventory.
	//	Create().
	//	SetItemID(1).
	//	SetQuantity(100).
	//	SetPrice(9.99).
	//	Save(context.Background())
	//if err != nil {
	//	log.Fatalf("failed creating inventory item: %v", err)
	//}
	//log.Printf("New Inventory Item: %+v\n", newItem)

}
