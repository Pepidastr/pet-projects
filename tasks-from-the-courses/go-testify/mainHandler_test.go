package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"strings"
	"testing"
)

func TestMainHandlerStatusOKNotEmptyBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	status := responseRecorder.Code
	bodyString := responseRecorder.Body.String()

	require.Equal(t, status, http.StatusOK)
	assert.NotEmpty(t, bodyString)
}

func TestMainHandWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=1&city=EKB", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	status := responseRecorder.Code
	expected := `wrong city value`
	assert.Equal(t, status, http.StatusBadRequest)
	assert.Equal(t, responseRecorder.Body.String(), expected)

}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := len(cafeList["moscow"])
	req := httptest.NewRequest("GET", "/cafe?count="+fmt.Sprint(totalCount+100)+"&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	bodyString := responseRecorder.Body.String()
	bodylist := strings.Split(bodyString, ",")

	assert.Equal(t, len(bodylist), totalCount)
}
