package main

import (
	"context"
	"log"
	"time"

	"hotels-api/clients/queues"
	controllersHotels "hotels-api/controllers/hotels"
	controllersReservations "hotels-api/controllers/reservations"
	repositoriesHotels "hotels-api/repositories/hotels"
	repositoriesReservations "hotels-api/repositories/reservations"
	servicesHotels "hotels-api/services/hotels"
	servicesReservations "hotels-api/services/reservations"

	"github.com/gin-contrib/cors" // Importa el paquete de CORS
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Conexión MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:root@mongo:27017"))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Configuración de Repositorios
	hotelsRepo := repositoriesHotels.NewMongo(repositoriesHotels.MongoConfig{
		Host:       "mongo",
		Port:       "27017",
		Username:   "root",
		Password:   "root",
		Database:   "hotels-api",
		Collection: "hotels",
	})

	reservationsRepo := repositoriesReservations.NewMongo(mongoClient, "hotels-api", "reservations")

	// Configuración de Cache y RabbitMQ
	cacheRepo := repositoriesHotels.NewCache(repositoriesHotels.CacheConfig{
		MaxSize:      100000,
		ItemsToPrune: 100,
		Duration:     30 * time.Second,
	})
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "root",
		Password:  "root",
		QueueName: "hotels-news",
	})

	// Servicios
	hotelsService := servicesHotels.NewService(hotelsRepo, cacheRepo, eventsQueue)
	reservationsService := servicesReservations.NewService(reservationsRepo)

	// Controladores
	hotelsController := controllersHotels.NewController(hotelsService)
	reservationsController := controllersReservations.NewController(reservationsService)

	// Rutas
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "*"}, // Permite localhost y cualquier origen
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas de Reservas y Hoteles (usando solo `hotel_id` en las rutas para evitar conflictos)
	router.POST("/reservations", reservationsController.CreateReservation)
	router.GET("/hotels/:hotel_id", hotelsController.GetHotelByID)
	router.POST("/hotels", hotelsController.Create)
	router.PUT("/hotels/:hotel_id", hotelsController.Update)
	router.DELETE("/hotels/:hotel_id", hotelsController.Delete)
	router.GET("/users/:user_id/reservations", reservationsController.GetReservationsByUserID)

	// Ejecutar servidor
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
