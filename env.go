//go:build windows

package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

//go:embed static/*
var static embed.FS

//go:generate npx tailwindcss build -i static/css/style.css -o static/css/styles.css --minify

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Init for .env done")
}