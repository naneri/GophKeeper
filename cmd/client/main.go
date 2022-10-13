package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/naneri/GophKeeper/cmd/client/apiClient"
	"github.com/naneri/GophKeeper/cmd/client/config"
	"github.com/naneri/GophKeeper/cmd/client/jobs"
	"github.com/naneri/GophKeeper/cmd/client/router"
	"github.com/naneri/GophKeeper/cmd/client/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	printBuildData()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	idleConnsClosed := make(chan struct{})

	var cfg config.Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalln("Error parsing config: ", err.Error())
	}

	appRouter := router.Router{}
	var server *http.Server

	recordStorage := storage.RecordStorage{}

	server = &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: appRouter.GetHandler(&cfg, &recordStorage),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		var client *apiClient.ApiClient
		for {
			recordErr := jobs.UpdateRecords(&cfg, &recordStorage, client)
			if recordErr != nil {
				log.Println("error updating records")
			}
			time.Sleep(3 * time.Second)
		}
	}()

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

func printBuildData() {
	fmt.Println("Build version: " + buildVersion)
	fmt.Println("Build date: " + buildDate)
}
