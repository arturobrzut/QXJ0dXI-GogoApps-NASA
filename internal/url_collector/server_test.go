package url_collector

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func getRouter() *gin.Engine {
	config, _ := Setup("DEMO_KEY", "8080", "5")
	r := router(config)
	return r
}

func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if !f(w) {
		t.Fail()
	}
}

func TestServNotFound(t *testing.T) {
	r := getRouter()

	req, _ := http.NewRequest("GET", "/somepath", nil)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusNotFound
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "page not found") > 0
		return statusOK && pageOK
	})
}

func TestServOk(t *testing.T) {
	r := getRouter()

	req, _ := http.NewRequest("GET", "/pictures?from=2021-01-01&to=2021-01-03", nil)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "urls") > 0

		return statusOK && pageOK
	})

}

func TestServWrongBadDate(t *testing.T) {
	r := getRouter()

	req, _ := http.NewRequest("GET", "/pictures?from=2021:01:01&to=2021:01:03", nil)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusBadRequest
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "parsing time") > 0
		return statusOK && pageOK
	})

}

func TestServWrongWrongToDate(t *testing.T) {
	r := getRouter()

	req, _ := http.NewRequest("GET", "/pictures?from=2021-01-01&to=2020-01-03", nil)
	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusBadRequest
		p, err := ioutil.ReadAll(w.Body)
		pageOK := err == nil && strings.Index(string(p), "Error:Field validation") > 0
		fmt.Println(string(p))
		return statusOK && pageOK
	})

}
