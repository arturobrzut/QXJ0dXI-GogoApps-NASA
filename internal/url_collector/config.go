package url_collector

import "strconv"

func Setup(apiKey, port, cr string) (*serviceConfig, error) {
	config := serviceConfig{}
	var err error
	config.apiKey = apiKey
	config.serverURL = "https://api.nasa.gov/planetary/apod"
	config.port, err = strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	config.concurrentRequest, err = strconv.Atoi(cr)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
