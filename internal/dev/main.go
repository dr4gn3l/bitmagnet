package main

import (
	// "github.com/bitmagnet-io/bitmagnet/internal/dev/app"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	db, _ := gorm.Open(postgres.Open("host=databases port=5432 user=bitmagnet_admin password=S31df1KcHAuUuv_Ki9 dbname=bitmagnet sslmode=disable"))
	var torrents []model.Torrent
	db.Limit(100).Find(torrents)
	for _, torrent := range torrents {
		fmt.Println(torrent.Name)
	}
	// app.New().Run()
}
