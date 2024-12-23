package main

import (
	"log"

	services "github.com/alphabatem/nft-proxy/service"
	"github.com/babilu-online/common/context"
	"github.com/joho/godotenv"
)

func LoadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	LoadEnvironment()

	ctx, err := context.NewCtx(
		&services.SqliteService{},
		&services.StatService{},
		&services.ResizeService{},
		&services.SolanaService{},
		&services.SolanaImageService{},
		&services.ImageService{},
		&services.HttpService{},
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = ctx.Run()
	log.Fatal(err)
}
