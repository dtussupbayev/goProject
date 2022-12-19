package main

import (
	"context"
	"github.com/Assyl00/goProject/internal/http"
	"github.com/Assyl00/goProject/internal/message_broker/kafka"
	"github.com/Assyl00/goProject/internal/store/postgres"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	//r := gin.Default()
	//r.POST("/login", handler.LoginHandler)
	//
	//r.GET("/api/product",
	//	middleware.ValidateToken(),
	//	middleware.Authorization([]int{4}),
	//	handler.GetAllProducts)
	//
	//r.Run("localhost:8080")

	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	urlExample := "postgres://postgres:postgres@localhost:5432/shop"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	//manager, err := auth.NewManager(key)
	//if err != nil {
	//	panic(err)
	//}

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer2")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	srv := http.NewServer(
		context.Background(),
		http.WithAddress(":8081"),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithBroker(broker),
	)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
	log.Println("Pinged DB")

}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}
