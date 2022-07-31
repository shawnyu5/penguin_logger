package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"server/utils"

	"github.com/gocolly/colly"
)

// TestLoggerHandler tests the logger handler, and checks if the product obj we get back is correct
func TestLoggerHandler(t *testing.T) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.penguinmagic.com", "www.penguinmagic.com/openbox/"),
	)

	penguinProduct := LoggerProduct{} // product on penguinmagic

	utils.GetTitle(c, &penguinProduct.Title)
	utils.GetPrice(c, &penguinProduct.Original_price)
	utils.GetDiscountedPrice(c, &penguinProduct.Discount_price)
	utils.GetDiscountPercentage(c, &penguinProduct.Discount_percentage)

	c.Visit("https://www.penguinmagic.com/openbox/")

	req := httptest.NewRequest(http.MethodGet, "/logger", nil)
	w := httptest.NewRecorder()
	loggerHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	var loggerProduct LoggerProduct
	err = json.Unmarshal(data, &loggerProduct) // parse json data into struct
	if err != nil {
		t.Errorf("Error unmarshalling json: %v", err)
	}

	if loggerProduct != penguinProduct {
		t.Errorf("Expected %v, got %v", loggerProduct, data)
	}
}

// TestSearchHandlerSuccess tests the search handler with a valid query
func TestSearchHandlerSuccess(t *testing.T) {
	// a successful search should return a list of products
	req := httptest.NewRequest(http.MethodPost, "/search", strings.NewReader("hello"))
	w := httptest.NewRecorder()
	searchHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// if we didnt get anything back, something has gone wrong
	if len(data) == 0 {
		t.Errorf("Expected data, got %v", data)
	}

}

// TestSearchHandlerFail tests the search handler with a invalid request method
func TestSearchHandlerFail(t *testing.T) {
	// a successful search should return a list of products
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	w := httptest.NewRecorder()
	searchHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %v", res.StatusCode)
	}
}
