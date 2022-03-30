package url_collector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetup(t *testing.T) {
	_, err := Setup("API_KEY", "8080", "5")
	assert.Equal(t, err, nil, "Error should be nil")

	_, err = Setup("API_KEY", "abc", "5")
	assert.NotEqual(t, err, nil, "port should be numerical value")

	_, err = Setup("API_KEY", "8080", "abc")
	assert.NotEqual(t, err, nil, "cr should be numerical value")
}
