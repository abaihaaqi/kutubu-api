# Cara menjalankan API

## Environment Configuration

1. Duplikat file .env.example
2. Rename menjadi .env
3. Sesuaikan konfigurasi dengan PostgreSQL yang terinstal

## Get Packages

```
go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/postgres github.com/joho/godotenv github.com/gin-contrib/cors
```

## Run App

```
go run main.go database.go model.go handler.go
```
