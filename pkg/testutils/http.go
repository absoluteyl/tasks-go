package testutils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetupAuthorizationHeader(req *http.Request, authStr string) {
	req.Header.Add("Authorization", authStr)
}

func HttpStatusShouldBe(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, rr.Code, "Handler returned wrong status code")
}

func HttpResponseShouldBe(t *testing.T, rr *httptest.ResponseRecorder, expectedBody string) {
	assert.Equal(t, expectedBody, rr.Body.String(), "Handler returned unexpected body")
}
