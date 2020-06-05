package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleWare_ServeHTTP(t *testing.T) {
	mockHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {})
	wrapper := &CORS{Handler: mockHandler}
	req := httptest.NewRequest(http.MethodOptions, "/v1", nil)
	rr := httptest.NewRecorder()
	wrapper.ServeHTTP(rr, req)
	resp := rr.Result()

	allowOriginHeader := resp.Header.Get("Access-Control-Allow-Origin")
	if allowOriginHeader != "*" {
		t.Errorf("allow origin header must be *")
	}

	allowMethodsHeader := resp.Header.Get("Access-Control-Allow-Methods")
	if allowMethodsHeader != "GET, PUT, POST, PATCH, DELETE" {
		t.Errorf("allow methods header must be \"GET, PUT, POST, PATCH, DELETE\"")
	}

	allowHeadersHeader := resp.Header.Get("Access-Control-Allow-Headers")
	if allowHeadersHeader != "Content-Type, Authorization" {
		t.Errorf("allow headers must be set to \"Content-Type, Authorization\"")
	}

	exposeHeadersHeader := resp.Header.Get("Access-Control-Expose-Headers")
	if exposeHeadersHeader != "Authorization" {
		t.Errorf("expose header must be \"Authorization\"")
	}

	maxAgeHeader := resp.Header.Get("Access-Control-Max-Age")
	if maxAgeHeader != "600" {
		t.Errorf("Max age for browser usage must be 600")
	}
}
