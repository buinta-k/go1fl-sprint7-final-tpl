package main

import (
	"fmt"
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

	city := "moscow"

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(len(cafeList[city]), 100)},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()

		url := fmt.Sprintf("/cafe?count=%d&city=%s", v.count, city)
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		bodyStr := response.Body.String()

		gotCount := len(strings.Split(bodyStr, ","))
		if bodyStr == "" {
			gotCount = 0
		}

		assert.Equal(t, v.want, gotCount)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		wantCount int
		word      string
	}{
		{0, "фасоль"},
		{2, "кофе"},
		{1, "вилка"},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()

		url := fmt.Sprintf("/cafe?city=moscow&search=%s", v.word)
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		bodyStr := strings.TrimSpace(response.Body.String())

		if v.wantCount == 0 {
			assert.Empty(t, bodyStr)
			continue
		}

		cafes := strings.Split(bodyStr, ",")

		assert.Len(t, cafes, v.wantCount)

		for _, name := range cafes {
			cafename := strings.ToLower(name)

			assert.Contains(t, cafename, v.word)
		}
	}
}
