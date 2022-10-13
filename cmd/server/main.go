package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
	"github.com/naneri/GophKeeper/cmd/server/config"
	"github.com/naneri/GophKeeper/cmd/server/router"
	"github.com/naneri/GophKeeper/internal/app/record"
	"github.com/naneri/GophKeeper/internal/app/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	idleConnsClosed := make(chan struct{})

	cfg, err := getConfig()
	if err != nil {
		log.Fatalln("Error parsing config: " + err.Error())
	}

	appRouter := router.Router{}
	var server *http.Server
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("error starting database: " + err.Error())
	}
	_ = db.AutoMigrate(&user.User{}, &record.Record{})

	server = &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: appRouter.GetHandler(db, cfg),
	}

	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if serverShutDownErr := server.Shutdown(context.Background()); serverShutDownErr != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown: %v", serverShutDownErr)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	if listenErr := server.ListenAndServe(); listenErr != http.ErrServerClosed {
		// ошибки старта или остановки Listener
		log.Fatalf("HTTP server ListenAndServe: %v", listenErr)
	}

	<-idleConnsClosed
	fmt.Println("Server Shutdown gracefully")
}

func getConfig() (*config.Config, error) {
	var cfg config.Config
	err := env.Parse(&cfg)

	if err != nil {
		return &cfg, err
	}

	log.Printf("Server starting with params: %+v", cfg)

	return &cfg, nil
}
