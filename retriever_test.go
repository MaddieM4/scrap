package scrap

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupHttpServer(t *testing.T, data []byte) *httptest.Server {
	var handler func(w http.ResponseWriter, r *http.Request)
	if len(data) > 0 {
		handler = func(w http.ResponseWriter, r *http.Request) {
			var sent int
			var total int = len(data)
			for sent < total {
				sent_this_round, err := w.Write(data[sent:])
				if err != nil {
					t.Fatal(err)
				}
				sent += sent_this_round
			}
		}
	} else {
		handler = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Four Oh Four", 404)
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	return ts
}

func serverAuth(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.Write([]byte("No auth"))
	} else {
		w.Write([]byte(auth))
	}
}

func TestHttpRetriever(t *testing.T) {
	ts := setupHttpServer(t, []byte(sample_html))
	defer ts.Close()

	trq := NewTestRQ()
	req := trq.CreateRequest(ts.URL)

	t.Logf("Requesting: %s", ts.URL)
	resp, err := HttpRetriever(req)
	if err != nil {
		t.Fatal(err)
	}

	got_response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	compare(t, []byte(sample_html), got_response)
}

func TestHttpRetriever_NoServer(t *testing.T) {
	url := "http://no.such.host/"
	trq := NewTestRQ()
	req := trq.CreateRequest(url)

	_, err := HttpRetriever(req)
	if err == nil {
		t.Fatal("Should have failed, didn't")
	}
}

func TestHttpRetriever_ErrorCode(t *testing.T) {
	ts := setupHttpServer(t, []byte{}) // Responds with 404
	defer ts.Close()

	url := ts.URL
	trq := NewTestRQ()
	req := trq.CreateRequest(url)

	t.Logf("Requesting: %s", url)
	resp, err := HttpRetriever(req)
	if err != nil {
		t.Fatal(err)
	}
	compare(t, 404, resp.StatusCode)
}

func TestHttpRetriever_Auth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(serverAuth))
	defer ts.Close()

	url := ts.URL
	trq := NewTestRQ()
	req := trq.CreateRequest(url)

	getContents := func(req ScraperRequest) string {
		resp, err := HttpRetriever(req)
		if err != nil {
			t.Fatal(err)
		}
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		return string(contents)
	}

	compare(t, "No auth", getContents(req))

	req.Auth = &RequestAuth{
		Username: "Flibberty",
		Password: "Jibbit",
	}
	compare(t, "Basic RmxpYmJlcnR5OkppYmJpdA==", getContents(req))
}
