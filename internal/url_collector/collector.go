package url_collector

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func collect(config *serviceConfig, pd picturesDate) ([]string, error) {
	var dates []string
	for value := pd.From; value.Before(pd.To); value = value.Add(time.Hour * 24) {
		dates = append(dates, value.Format("2006-01-02"))
	}
	dates = append(dates, pd.To.Format("2006-01-02"))
	return downloadAllData(config, &dates)
}

func downloadAllData(config *serviceConfig, dateList *[]string) ([]string, error) {
	dateChan := make(chan string, len(*dateList))
	resultsChan := make(chan result, len(*dateList))
	var urls, errs []string

	for w := 1; w <= config.concurrentRequest; w++ {
		go func(configuration *serviceConfig, dates <-chan string, results chan<- result) {
			for date := range dates {
				res := result{}
				res.url, res.err = downloadData(configuration, date)
				results <- res
			}
		}(config, dateChan, resultsChan)
	}
	for _, date := range *dateList {
		dateChan <- date
	}
	close(dateChan)

	for range *dateList {
		res := <-resultsChan
		if res.url != "" {
			urls = append(urls, res.url)
		}
		if res.err != nil {
			errs = append(errs, res.err.Error())
		}
	}
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, ", "))
	}
	return urls, nil
}

func downloadData(config *serviceConfig, date string) (string, error) {
	request, err := http.NewRequest("GET", config.serverURL, nil)
	if err != nil {
		return "", err
	}
	query := request.URL.Query()
	query.Add("api_key", config.apiKey)
	query.Add("date", date)
	request.URL.RawQuery = query.Encode()

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New("API Error: " + response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", errors.New("API Error: " + string(body))
	}

	jsonData := map[string]string{}
	err = json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		return "", err
	}

	return jsonData["url"], nil
}
