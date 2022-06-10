package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DbProduct struct {
	Title            string  `json:"title"`
	Appearances      int32   `json:"appearances"`
	Average_discount float64 `json:"average_discount"`
	Average_price    float64 `json:"average_price"`
	Created_at       string  `json:"created_date"`
	Updated_at       string  `json:"updated_date"`
}

type PenguinProduct struct {
	Title              string
	Description        string
	OriginalPrice      float64
	DiscountPrice      float64
	DiscountPercentage float64
	Rating             int64
	IsValid            bool
	Reason             string
}

// Connection URI
var uri string

func main() {
	client := connectDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// var p Product
	var penguinProduct PenguinProduct
	p := getProduct()
	if err := json.Unmarshal(p, &penguinProduct); err != nil {
		panic(err)
	}
	// fmt.Println(fmt.Sprintf("main penguinProduct: %v", penguinProduct)) // __AUTO_GENERATED_PRINT_VAR__

	coll := client.Database("penguin_magic").Collection("open_box")
	var dbResult bson.M
	if err := coll.FindOne(context.TODO(), bson.D{{"title", penguinProduct.Title}}).Decode(&dbResult); err != nil {
		const appearances int32 = 0
		const average_discount float64 = 0
		const average_price float64 = 0
		var created_at string = time.Now().Format("2006-01-02 15:04:05")
		var updated_at string = time.Now().Format("2006-01-02 15:04:05")

		dbResult = bson.M{"title": penguinProduct.Title, "appearances": appearances, "average_discount": average_discount, "average_price": average_price, "created_at": created_at, "updated_at": updated_at}

	}
	dbProduct := constructProductObj(dbResult)
	fmt.Println(fmt.Sprintf("main dbProduct: %+v", dbProduct)) // __AUTO_GENERATED_PRINT_VAR__
	err := updateProduct(&dbProduct, penguinProduct)
	if err != nil {
		panic(err)
	}
	return
	saveProduct(&dbProduct, coll)

}

// getProduct will get the current product from penguin magic
func getProduct() []byte {
	// make http request to get product
	// TODO: consider making custom end point for logger
	res, err := http.Get(os.Getenv("API_URL") + "/coinProduct")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)

	return body
}

// constructProductObj construct a product object from the mongodb result
func constructProductObj(b bson.M) DbProduct {
	product := DbProduct{Title: b["title"].(string), Appearances: b["appearances"].(int32), Average_discount: b["average_discount"].(float64), Average_price: b["average_price"].(float64)}

	// not all products has these time stamps
	if b["created_at"] == nil {
		product.Created_at = time.Now().Format("2006-01-02 15:04:05")
	} else {
		product.Created_at = b["created_at"].(string)
	}

	if b["updated_at"] == nil {
		product.Updated_at = time.Now().Format("2006-01-02 15:04:05")
	} else {
		product.Updated_at = b["updated_at"].(string)

	}
	return product
}

func connectDB() *mongo.Client {
	godotenv.Load()
	uri = os.Getenv("MONGODB_URI")
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return client
}

// saveProduct save the dbproduct to the collection
func saveProduct(dbProduct *DbProduct, coll *mongo.Collection) {
	b, err := bson.Marshal(dbProduct)
	if err != nil {
		panic(err)
	}

	filter := bson.D{{"title", dbProduct.Title}}
	result, err := coll.ReplaceOne(context.TODO(), filter, b)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("saveProduct result: %v", result.UpsertedID)) // __AUTO_GENERATED_PRINT_VAR__

}

// updateProduct update the dbproduct with the new product passed in
func updateProduct(dbProduct *DbProduct, penguinProduct PenguinProduct) error {
	if dbProduct.Title != penguinProduct.Title {
		// return an error
		return fmt.Errorf("Product titles do not match")
	} else if !penguinProduct.IsValid {
		return fmt.Errorf("Product is not valid: %v", penguinProduct.Reason)
	}
	// update the dbproduct with the new product
	dbProduct.Appearances = dbProduct.Appearances + 1
	dbProduct.Average_discount = (dbProduct.Average_discount*float64(dbProduct.Appearances-1) + penguinProduct.DiscountPrice) / float64(dbProduct.Appearances)
	dbProduct.Average_price = (dbProduct.Average_price*float64(dbProduct.Appearances-1) + penguinProduct.DiscountPrice) / float64(dbProduct.Appearances)
	// get current date and time
	dbProduct.Updated_at = time.Now().Format("2006-01-02 15:04:05")

	return nil
}