package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	check_coin_product "server/coin_products"

	"server/utils"

	"server/search"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
)

var storage *cache.Cache

type LoggerProduct struct {
	Title               string  `json:"title"`
	Original_price      float64 `json:"original_price"`
	Discount_price      float64 `json:"discount_price"`
	Discount_percentage float64 `json:"discount_percentage"`
}

type CoinProductService interface {
}

func main() {
	// initialize the cache
	storage = cache.New(cache.NoExpiration, 10*time.Minute)
	routes := make(map[string]func(http.ResponseWriter, *http.Request))
	routes["/"] = homeHandler(routes)
	routes["/coinProduct"] = coinProductHandler
	routes["/logger"] = loggerHandler
	routes["/search"] = searchHandler
	routes["/favicon.ico"] = doNothing
	for k, v := range routes {
		http.HandleFunc(k, v)
	}

	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("(server) Error loading .env file")
	}
	// get from env
	port := ":" + os.Getenv("PORT")
	// set default port to 8080
	if port == ":" {
		port = ":8080"
	}
	fmt.Println("LISTENING ON PORT " + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

// coinProductHandler is the handler for the /coinProduct endpoint
func coinProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 404)
		return
	}
	productInfo := check_coin_product.Check("https://www.penguinmagic.com/openbox/")
	// check if product has changed
	var title interface{}
	found := false

	if storage != nil {
		title, found = storage.Get("coin_product_title")
	}
	if found && title.(string) == productInfo.Title {
		productInfo.IsValid = false
		productInfo.Reason = "Product has not changed"
	}
	if storage != nil {
		storage.Set("coin_product_title", productInfo.Title, cache.DefaultExpiration)
	}
	log.Println("/coinProduct:", productInfo.Title)
	j, err := json.MarshalIndent(productInfo, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(j))
}

// doNothing is a do nothing function
func doNothing(w http.ResponseWriter, r *http.Request) {}

func homeHandler(routes map[string]func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		var list string
		for k := range routes {
			list += k + "\n"
		}
		log.Println("/:", list)
		fmt.Fprintln(w, list)
	}

}

func loggerHandler(w http.ResponseWriter, r *http.Request) {

	c := colly.NewCollector(
		colly.AllowedDomains("www.penguinmagic.com", "www.penguinmagic.com/openbox/", "https://www.penguinmagic.com/p/12449"),
	)
	product := LoggerProduct{}

	utils.GetTitle(c, &product.Title)
	utils.GetPrice(c, &product.Original_price)
	utils.GetDiscountedPrice(c, &product.Discount_price)
	utils.GetDiscountPercentage(c, &product.Discount_percentage)

	// c.Visit("https://www.penguinmagic.com/p/17318")
	// c.Visit("https://www.penguinmagic.com/p/12449")
	c.Visit("https://www.penguinmagic.com/openbox/")

	if storage != nil {
		storage.Set("product_title", product.Title, cache.DefaultExpiration)
	}
	j, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Println(string(j))
	fmt.Fprintln(w, string(j))
}

// searchHandler is the handler for the /search endpoint
// Returns a list of products that match the search query in json
func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method "+r.Method+" not allowed", 404)
		return
	}

	logger := log.New(os.Stdout, "", log.LUTC)

	var s search.SearchService
	s = search.SearchServiceImpl{}
	s = search.LoggingMiddleware{Logger: logger, Next: s}

	body, err := ioutil.ReadAll(r.Body)
	product := search.Product{Title: string(body)}
	result, err := s.SearchByRegex(&product)
	// result := search.SearchByRegex(&product)

	j, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, string(j))
}
