package main

import (
	"fmt"
	"log"
	"os"
	"user_management_system/database"
	"user_management_system/notifier"
	routes "user_management_system/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
func main(){
	err:=godotenv.Load(".env")
	if err!=nil{
		log.Fatal("Error loading in .env file")
	}
	port:=os.Getenv("PORT")
	if port==""{
		port="8000"
	}

	// 2. Load config
	cfg := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL"),
	}
	err= database.NewConnection(cfg)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	
	router:=gin.New()
	router.Use(gin.Logger())
	
	n := notifier.NewNotifier(100)
	v1 := router.Group("/api/v1")

	routes.AuthRoutes(v1,n)
	routes.UserRoutes(v1)
	
	router.GET("/health",func(c *gin.Context){
		c.JSON(200,gin.H{"Status":"ok"})
	})

	fmt.Printf("%s",fmt.Sprintf("Server is running from port no %s",port))
	router.Run(":"+port)
}