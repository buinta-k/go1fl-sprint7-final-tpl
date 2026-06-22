package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		count string
		want  int
	}{
		{"/cafe?city=tula&count=0", 0},
		{"/cafe?city=tula&count=1", 1},
		{"/cafe?city=moscow&count=2", 2},
		{"/cafe?city=moscow&count=100", 100},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.count, nil)

		handler.ServeHTTP(response, req)

		bodyStr := response.Body.String()

		var gotCount int
		if bodyStr != "" {
			slice := strings.Split(bodyStr, ",")
			gotCount = len(slice)
		} else {
			gotCount = 0
		}

		assert.Equal(t, v.want, gotCount)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
		word      string
	}{
		{"/cafe?city=moscow&search=фасоль", 0, "фасоль"},
		{"/cafe?city=moscow&search=кофе", 2, "кофе"},
		{"/cafe?city=moscow&search=вилка", 1, "вилка"},
	}

	for _, v := range requests {

		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.search, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		bodyStr := strings.TrimSpace(response.Body.String())

		if v.wantCount == 0 {
			assert.Empty(t, bodyStr)
			continue
		}

		cafes := strings.Split(bodyStr, ",")

		assert.Equal(t, v.wantCount, len(cafes))

		for _, name := range cafes {
			cafename := strings.ToLower(name)

			assert.True(t, strings.Contains(cafename, v.word))
		}

	}

}

