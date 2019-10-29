package main

import (
	"log"

	pg "github.com/go-pg/pg/v9"
	"github.com/mateeullahmalik/serviceMonitor/httpd/handler"
	"github.com/mateeullahmalik/serviceMonitor/httpd/monitor"
	"github.com/mateeullahmalik/serviceMonitor/httpd/persistence/storage"
)

func main() {
	services := []string{
		"https://www.google.com",
		"https://www.hotmail.com",
		"https://www.facebook.com",
		"https://www.twitter.com",
		"https://www.yahoo.com",
		"https://www.amazon.com",
		"https://www.wikipedia.com",
		"https://www.live.com",
		"https://www.youtube.com",
		"https://www.instagram.com",
		"https://www.baidu.com",
		"https://www.stackoverflow.com",
	}

	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "",
		Addr:     "localhost:5432",
	})

	if db == nil {
		log.Printf("Failed to connect with database")
	}

	defer db.Close()
	log.Printf("Connection to database successfull")

	storageHandler := storage.StorageHandler{
		Db: db,
	}

	if err := storageHandler.CreateServicePingsTable(); err != nil {
		log.Panic(err)
	}

	monitor.MonitorServices(&storageHandler, services...)

	router := handler.SetupRouter(&storageHandler)
	router.Run(":1010")
}
