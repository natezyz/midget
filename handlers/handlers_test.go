package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/natezyz/midget/storage"
)

func TestCorrectMethod(t *testing.T) {
	s := &storage.Map{}
	s.Init()

	var json = []byte(`{"url":"www.google.com"}`)
	req, err := http.NewRequest("GET", "/encode", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Encode(s)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":404,"message":"method not found"}`
	if rr.Body.String() != expected {
		t.Errorf("handler return wrong response: got %v want %v", rr.Body.String(), expected)
	}
}

func TestEncodePOST(t *testing.T) {
	s := &storage.Map{}
	s.Init()

	var json = []byte(`{"url":"www.google.com"}`)
	req, err := http.NewRequest("POST", "/encode", bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Encode(s)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"7299966772734287354"}`
	if rr.Body.String() != expected {
		t.Errorf("handler return wrong response: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDecode(t *testing.T) {
	s := &storage.Map{}
	s.Init()
	s.Store("www.google.com")

	req, err := http.NewRequest("GET", "/decode/7299966772734287354", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Decode(s)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":200,"message":"www.google.com"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned wrong response: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDecodeCodeNotFound(t *testing.T) {
	s := &storage.Map{}
	s.Init()
	s.Store("www.google.com")

	req, err := http.NewRequest("GET", "/decode/1299966772734287354", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Decode(s)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	expected := `{"status":204,"message":"invalid code 1299966772734287354"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned wrong response: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRedirect(t *testing.T) {
	s := &storage.Map{}
	s.Init()
	s.Store("www.google.com")

	req, err := http.NewRequest("GET", "/7299966772734287354", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Redirect(s)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMovedPermanently {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMovedPermanently)
	}
}
