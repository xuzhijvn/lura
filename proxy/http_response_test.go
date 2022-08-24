// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"compress/gzip"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luraproject/lura/v2/encoding"
)

func TestNopHTTPResponseParser(t *testing.T) {
	w := httptest.NewRecorder()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("header1", "value1")
		w.Write([]byte("some nice, interesting and long content"))
	}
	req, _ := http.NewRequest("GET", "/url", nil)
	handler(w, req)
	result, err := NoOpHTTPResponseParser(context.Background(), w.Result())
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !result.IsComplete {
		t.Error("unexpected result")
	}
	if len(result.Data) != 0 {
		t.Error("unexpected result")
	}
	if result.Metadata.StatusCode != http.StatusOK {
		t.Error("unexpected result")
	}
	headers := result.Metadata.Headers
	if h, ok := headers["Header1"]; !ok || h[0] != "value1" {
		t.Error("unexpected result:", result.Metadata.Headers)
	}
	body, err := ioutil.ReadAll(result.Io)
	if err != nil {
		t.Error("unexpected error:", err.Error())
	}
	if string(body) != "some nice, interesting and long content" {
		t.Error("unexpected result")
	}
}

func TestDefaultHTTPResponseParser_gzipped(t *testing.T) {
	w := httptest.NewRecorder()
	handler := func(w http.ResponseWriter, r *http.Request) {
		gzipWriter, _ := gzip.NewWriterLevel(w, gzip.BestSpeed)
		defer gzipWriter.Close()

		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		gzipWriter.Write([]byte(`{"msg":"some nice, interesting and long content"}`))
		gzipWriter.Flush()
	}
	req, _ := http.NewRequest("GET", "/url", nil)
	req.Header.Add("Accept-Encoding", "gzip")
	handler(w, req)

	result, err := DefaultHTTPResponseParserFactory(HTTPResponseParserConfig{
		Decoder:         encoding.JSONDecoder,
		EntityFormatter: DefaultHTTPResponseParserConfig.EntityFormatter,
	})(context.Background(), w.Result())

	if err != nil {
		t.Error(err)
	}

	if !result.IsComplete {
		t.Error("unexpected result")
	}
	if len(result.Data) != 1 {
		t.Error("unexpected result")
	}
	if m, ok := result.Data["msg"]; !ok || m != "some nice, interesting and long content" {
		t.Error("unexpected result")
	}
}

func TestDefaultHTTPResponseParser_plain(t *testing.T) {
	w := httptest.NewRecorder()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"msg":"some nice, interesting and long content"}`))
	}
	req, _ := http.NewRequest("GET", "/url", nil)
	handler(w, req)

	result, err := DefaultHTTPResponseParserFactory(HTTPResponseParserConfig{
		Decoder:         encoding.JSONDecoder,
		EntityFormatter: DefaultHTTPResponseParserConfig.EntityFormatter,
	})(context.Background(), w.Result())

	if err != nil {
		t.Error(err)
	}

	if !result.IsComplete {
		t.Error("unexpected result")
	}
	if len(result.Data) != 1 {
		t.Error("unexpected result")
	}
	if m, ok := result.Data["msg"]; !ok || m != "some nice, interesting and long content" {
		t.Error("unexpected result")
	}
}
