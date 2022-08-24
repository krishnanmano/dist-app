package main

import (
	"context"
	"dist-app/controllers"
	memlist "dist-app/memberlist"
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

var router *gin.Engine

const (
	GracefulShutdownTimeout = 5 * time.Second
	RequestHeaderTimeout    = 500 * time.Millisecond
)

func main() {
	var appPort string
	flag.StringVar(&appPort, "app-port", "8080", "port to initiate app instance")

	var gossipPort int
	flag.IntVar(&gossipPort, "gossip-port", 7950, "port to initiate app instance") //nolint

	var gossipLeader string
	flag.StringVar(&gossipLeader, "gossip-leader", "", "leader service 'ip:port'")
	flag.Parse()

	gossipNode := memlist.NewGossipNode(gossipPort, gossipLeader, appPort)

	distAppService := service.NewDistAppService()
	distAppController := controllers.NewDistappController(distAppService, gossipNode)
	gossipController := controllers.NewGossipController(distAppService)

	router = gin.Default()
	mapurls(distAppController, gossipController)

	server := &http.Server{
		Addr:              ":" + appPort,
		Handler:           router,
		ReadHeaderTimeout: RequestHeaderTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	gracefulShutdown(server, gossipNode)
}

func mapurls(distAppController *controllers.DistappController, gossipController *controllers.GossipController) {
	serviceGroup := router.Group("/dist-app-service/api")
	serviceGroup.GET("/health", distAppController.Health)
	serviceGroup.GET("/messages", distAppController.GetMessages)
	serviceGroup.POST("/messages", distAppController.SaveMessage)

	gossipGroup := router.Group("/gossip/api")
	gossipGroup.GET("/health", gossipController.Health)
	gossipGroup.POST("/publish", gossipController.PublishMessage)
}

func gracefulShutdown(server *http.Server, gossipNode memlist.GossipNode) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //nolint
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")

	gossipNode.GracefullyLeave(GracefulShutdownTimeout)
}
