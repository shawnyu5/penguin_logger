package search

import (
	"context"
	"log"
	"os"
	"testing"
)

func beforeEach() SearchService {
	logger := log.New(os.Stdout, "", log.LUTC)

	var ss SearchService
	ss = SearchServiceImpl{}
	ss = LoggingMiddleware{Logger: logger, Next: ss}
	return ss
}

// TestAbleToConnectToDb tests if we can connect to the database
func TestAbleToConnectToDb(t *testing.T) {
	ss := beforeEach()
	client := ss.connectDB()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if client == nil {
		t.Error("Failed to connect to mongodb")
	}
}

func TestSearchByRegex(t *testing.T) {
	ss := beforeEach()
	product := &Product{Title: "card"}
	found, err := ss.SearchByRegex(product)
	if err != nil {
		t.Error("Error searching for product")
	}

	if len(found) == 0 {
		t.Error("Failed to find product")
	}
}
