# Athletics backend
Backend application for athletics app, written in Golang

## Setup
1. download dependencies `go mod download`
2. install `docker` and `docker compose`
3. create 
   - `.env` file from `.env.sample`
   - `.env.test` file from `.env.test.sample`

## Development
1. run postgres database  
   `docker compose up db`
2. run the app (call in project root directory)  
   `go run ./cmd`

Alternatively, run the app with the database in docker (this simulates production deployment):  
`docker compose up`

## Testing
`go test ./...`

## App structure
Because appplication is written in golang, the easiest way to explore the app is to start from `./cmd/main.go` file and follow along to check how things work