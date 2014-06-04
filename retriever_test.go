package scrap

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BrokenReader struct{}

func (br BrokenReader) Read([]byte) (int, error) {
	return 0, errors.New("BrokenReader says hello")
}

func TestParseReader(t *testing.T) {
	// Success case is verified by other tests, just test error
	_, err := parseReader(ScraperRequest{}, BrokenReader{})
	if err == nil {
		t.Fatal("Should not have suceeded")
	}
	compare(t, "BrokenReader says hello", err.Error())
}

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

func TestHttpRetriever(t *testing.T) {
	ts := setupHttpServer(t, []byte(sample_html))
	defer ts.Close()

	trq := NewTestRQ()
	req := trq.CreateRequest(ts.URL)

	t.Logf("Requesting: %s", ts.URL)
	n, err := HttpRetriever(req)
	if err != nil {
		t.Fatal(err)
	}
	compare(t, 3, len(n.Find("a")))
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
	_, err := HttpRetriever(req)
	if err == nil {
		t.Fatal("Should have failed, didn't")
	}
}
