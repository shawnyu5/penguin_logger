package middware

import (
	"log"
	"net/http"
	"time"
)

// a logger middleware
type LoggerMidware struct {
	log  log.Logger
	next func(w http.ResponseWriter, r *http.Request)
}

func (mw LoggerMidware) logger(method string, responseCode string, err error, responseBody string) {
	defer func(begin time.Time) {
		log.Printf("method=%s, responseCode=%s, err=%s, responseBody=%s, time=%s", method, responseCode, err, responseBody, begin)
	}(time.Now())
}
