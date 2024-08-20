package main

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/storeCustomer/config"
	handlers "github.com/tanush-128/openzo_backend/storeCustomer/internal/api"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/middlewares"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/pb"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/repository"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/service"
	"google.golang.org/grpc"
)

var UserClient pb.UserServiceClient

type User2 struct {
}

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	db, err := connectToDB(cfg) // Implement database connection logic
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}

	// // Initialize gRPC server
	// grpcServer := grpc.NewServer()
	// StoreCustomerpb.RegisterStoreCustomerServiceServer(grpcServer, service.NewGrpcStoreCustomerService(StoreCustomerRepository, StoreCustomerService))
	// reflection.Register(grpcServer) // Optional for server reflection

	conf := ReadConfig()
	p, _ := kafka.NewProducer(&conf)

	// go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	//Initialize gRPC client
	conn, err := grpc.Dial(cfg.UserGrpc, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)
	UserClient = c

	storeCustomerRepository := repository.NewStoreCustomerRepository(db)
	StoreCustomerService := service.NewStoreCustomerService(storeCustomerRepository, p)

	// Initialize HTTP server with Gin
	router := gin.Default()
	handler := handlers.NewHandler(&StoreCustomerService)

	router.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// router.Use(middlewares.JwtMiddleware(c))

	router.GET("/store/:id", handler.GetStoreCustomersByStoreID)
	router.POST("/", handler.CreateStoreCustomer)
	router.GET("/:id", handler.GetStoreCustomerByID)
	router.PUT("/", handler.UpdateStoreCustomer)
	router.Use(middlewares.NewMiddleware(c).JwtMiddleware)
	router.GET("/user", handler.GetStoreCustomersByUserID)
	router.DELETE("/:id", handler.DeleteStoreCustomer)

	// router.Use(middlewares.JwtMiddleware)

	router.Run(fmt.Sprintf(":%s", cfg.HTTPPort))

}
