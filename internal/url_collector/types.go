package url_collector

import "time"

type serviceConfig struct {
	port              int
	apiKey            string
	concurrentRequest int
	serverURL         string
}

type picturesDate struct {
	From time.Time `form:"from" binding:"required" time_format:"2006-01-02"`
	To   time.Time `form:"to" binding:"required,gtefield=From" time_format:"2006-01-02"`
}

type result struct {
	url string
	err error
}
