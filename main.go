package main

import (
	"context"
	"dist-app/controllers"
	"dist-app/db/cockroachdb"
	"dist-app/logger"
	memlist "dist-app/memberlist"
	"dist-app/middleware"
	"dist-app/model"
	"dist-app/service"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	GracefulShutdownTimeout = 5 * time.Second
	RequestHeaderTimeout    = 500 * time.Millisecond
)

var (
	router *gin.Engine
)

func main() {
	var appPort string
	flag.StringVar(&appPort, "app-port", "8080", "port to initiate app instance")

	var gossipPort int
	flag.IntVar(&gossipPort, "gossip-port", 7950, "port to initiate app instance") //nolint

	var gossipLeader string
	flag.StringVar(&gossipLeader, "gossip-leader", "", "leader service 'ip:port'")
	flag.Parse()

	// Initialize logger
	logger.InitDefaultLogger("app.log", true, map[string]interface{}{
		"service": "dist-app",
	})

	// Database Layer
	cockroachdb.InitCockroachDB()
	transactionDao := cockroachdb.NewTransactionDAO()
	cockroachdb.DBClient.AutoMigrate(&model.Transaction{})

	// Gossip Layer
	//gossipNode := memlist.NewGossipNode(gossipPort, gossipLeader, appPort)

	//Service Layer
	txnService := service.NewTransactionService(nil, transactionDao)

	//Controller Layer
	txnController := controllers.NewTransactionController(txnService)
	//gossipController := controllers.NewGossipController(txnService)

	// Gin configuration
	//if _, ok := os.LookupEnv("GIN_MODE"); !ok {
	//	gin.SetMode(gin.ReleaseMode)
	//}
	router = gin.New()
	router.Use(middleware.JSONLogMiddleware(), gin.Recovery())
	mapurls(nil, txnController)

	server := &http.Server{
		Addr:              ":" + appPort,
		Handler:           router,
		ReadHeaderTimeout: RequestHeaderTimeout,
		ReadTimeout:       RequestHeaderTimeout,
		WriteTimeout:      RequestHeaderTimeout,
		IdleTimeout:       RequestHeaderTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	gracefulShutdown(server, nil)
}

func mapurls(gossipController *controllers.GossipController, txnController *controllers.TransactionController) {
	transactionServiceGroup := router.Group("/transaction-service/api")
	{
		transactionServiceGroup.GET("/health", txnController.Health)
		transactionServiceGroup.GET("/transactions", txnController.GetTransactions)
		transactionServiceGroup.GET("/transactions/:id", txnController.FindTransactionByID)
		transactionServiceGroup.POST("/transactions", txnController.SaveTransaction)
		transactionServiceGroup.DELETE("/transactions/:id", txnController.DeleteTransactionByID)
		transactionServiceGroup.PUT("/transactions/:id", txnController.UpdateTransactionByID)
		transactionServiceGroup.PATCH("/transactions/:id", txnController.UpdateTransactionByID)
	}

	if gossipController != nil {
		gossipGroup := router.Group("/gossip/api")
		{
			gossipGroup.GET("/health", gossipController.Health)
			gossipGroup.POST("/publish", gossipController.PublishMessage)
		}
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("*********** logger1", time.Now())
	}
}

func Logger2() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("*********** logger2", time.Now())
	}
}

func gracefulShutdown(server *http.Server, gossipNode *memlist.GossipNode) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //nolint
	<-quit

	log.Println("Database service: shutdown initiated...")
	db, _ := cockroachdb.DBClient.DB()
	if err := db.Close(); err != nil {
		log.Println("Database Service: forcing shutdown, ", err)
	}
	log.Println("Database service: shutdown completed...")

	log.Println("Application service: shutdown initiated...")
	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Application Service: forcing shutdown, ", err)
	}
	log.Println("Application service: shutdown completed...")

	if gossipNode != nil {
		log.Println("Gossip service: shutdown initiated...")
		gossipNode.GracefullyLeave(GracefulShutdownTimeout)
		log.Println("Gossip service: shutdown completed...")
	}
}
