package search

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type SearchService interface {
	// connectDB will connect to the mongodb
	connectDB() *mongo.Client
	// SearchByRegex will search for a product title by a regex
	// Returns an array of products
	SearchByRegex(*Product) ([]Product, error)
}

type SearchServiceImpl struct{}

type Product struct {
	Title               string  `json:"title"`
	Average_discount    float64 `json:"average_discount"`
	Average_price       float64 `json:"average_price"`
	Discount_percentage float64 `json:"discount_percentage"`
	Appearances         int32   `json:"appearances"`
}

// connectDB will connect to the mongodb
func (SearchServiceImpl) connectDB() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		panic("MONGODB_URI is not set")
	}
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	// log.Println("Database successfully connected and pinged.")
	return client
}

// SearchByRegex will search for a product title by a regex
// Returns an array of products
func (ss SearchServiceImpl) SearchByRegex(product *Product) ([]Product, error) {
	client := ss.connectDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		// log.Println("Database successfully disconnected.")
	}()

	coll := client.Database("penguin_magic").Collection("open_box")
	filter := bson.M{}

	filter["title"] = bson.M{"$regex": fmt.Sprintf(".*%s.*", product.Title)}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var products []Product
	// go through all products and add to array of products
	for cursor.Next(context.TODO()) {
		var singleProduct Product
		err := cursor.Decode(&singleProduct)
		if err != nil {
			return nil, err
		}
		products = append(products, singleProduct)
	}

	// fmt.Println(fmt.Sprintf("SearchByRegex products: %v", products)) // __AUTO_GENERATED_PRINT_VAR__
	return products, nil
}
