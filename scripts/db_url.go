package main

import (
	"encoding/json"
	"fmt"
	"os"

	"flag"

	"github.com/joho/godotenv"
)

// as atlasgo is incapable of reading .env files directly, this small script was added to make it work
func main() {
	envType := flag.String("env", "dev", "type of env: test, dev, prod")
	flag.Parse()

	sslMode := "disable"
	if *envType == "test" || *envType == "dev" {
		envFilePath := ".env"
		if *envType == "test" {
			envFilePath = ".env.test"
		}
		err := godotenv.Load(envFilePath)
		if err != nil {
			fmt.Println("Error loading .env file", err) // println allowed here
			os.Exit(1)
		}
	}

	envMap := map[string]string{
		"db_url": fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=public", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), sslMode),
	}
	w := os.Stdout
	if err := json.NewEncoder(w).Encode(&envMap); err != nil {
		os.Exit(1)
	}
}
