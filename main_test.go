package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)


requests := []struct {
		req   string
        want  int   // ожидаемое количество кафе в ответе
    }{
		{"/cafe?count=0&city=moscow", 0},
		{"/cafe?count=1&city=tula", 1},

		{"/cafe?count=2&city=moscow", 2},
		{"/cafe?count=100&city=tula", len(cafeList["tula"])},
    } 
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.req, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)

		var have int
		if response.Body.String() == "" {
			have = 0
		}else {
			have = len(strings.Split(response.Body.String(), ","))
		}
		assert.Equal(t, v.want, have)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		req    string
		search string 
		want   int 
	}{
		{"/cafe?city=moscow&search=фасоль", "фасоль", 0},
		{"/cafe?city=moscow&search=кофе", "кофе", 2},
		{"/cafe?city=moscow&search=вилка", "вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.req, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)

		have := 0
		body := strings.Split(response.Body.String(), ",")
		assert.Len(t, body, v.want)
		for _, value := range body {
			if assert.Contains(t, strings.ToLower(value), v.search) {
				have += 1
			}
		}
		assert.Equal(t, v.want, have)
	}
}