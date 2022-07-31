package search

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type LoggingMiddleware struct {
	Logger *log.Logger
	Next   SearchService
}

func (lm LoggingMiddleware) SearchByRegex(p *Product) (output []Product, err error) {
	defer func(begin time.Time) {
		lm.Logger.Printf("method=SearchByRegex input=%s output=%+v err=%s took=%s \n", p.Title, output, err, time.Since(begin))
	}(time.Now())

	output, err = lm.Next.SearchByRegex(p)
	return output, err
}

func (lm LoggingMiddleware) connectDB() *mongo.Client {
	lm.Logger.Println("Connecting to database")
	defer func() {
		lm.Logger.Println("Disconnecting from database")
	}()
	return lm.Next.connectDB()
}
