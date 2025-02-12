# Athletics backend
Backend application for athletics app, written in Golang

## Setup
1. download dependencies `go mod tidy`
2. install `docker` and `docker compose`
3. create 
   - `.env` file from `.env.sample`
   - `.env.test` file from `.env.test.sample`
4. install `atlasgo` to execute migrations
```shell
curl -sSf https://atlasgo.sh | sh
```

## Development
1. run postgres database  
   `docker compose up db`

2. execute migrations for dev environment
```bash
atlas migrate apply --env dev 
```

for testing
3. run the app (call in project root directory)  
   `go run ./cmd`

## Testing
Run migrations for test environment
```bash
atlas migrate apply --env test
```
`go test ./...`

## App structure
Because appplication is written in golang, the easiest way to explore the app is to start from `./cmd/main.go` file and follow along to check how things work