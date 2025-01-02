package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoStr := "mongodb://admin:admin@localhost:27017/routes?authSource=admin"
	println("Hello, World!")
}

