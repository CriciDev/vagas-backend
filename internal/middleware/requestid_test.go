package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDGeneratesWhenHeaderMissing(t *testing.T) {
	router := newRequestIDRouter()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(recorder, request)

	id := recorder.Header().Get(HeaderRequestID)
	if !validRequestID.MatchString(id) {
		t.Fatalf("expected a generated request id, got %q", id)
	}
}

func TestRequestIDReusesValidClientHeader(t *testing.T) {
	router := newRequestIDRouter()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	request.Header.Set(HeaderRequestID, "client-123")
	router.ServeHTTP(recorder, request)

	if got := recorder.Header().Get(HeaderRequestID); got != "client-123" {
		t.Fatalf("expected client id to be reused, got %q", got)
	}
}

func TestRequestIDReplacesInvalidClientHeader(t *testing.T) {
	router := newRequestIDRouter()

	invalid := "has spaces and symbols #!@"
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	request.Header.Set(HeaderRequestID, invalid)
	router.ServeHTTP(recorder, request)

	got := recorder.Header().Get(HeaderRequestID)
	if got == invalid {
		t.Fatal("expected invalid client id to be replaced")
	}
	if !validRequestID.MatchString(got) {
		t.Fatalf("expected a generated request id, got %q", got)
	}
}

func TestRequestIDHeaderBoundaries(t *testing.T) {
	cases := []struct {
		name   string
		header string
		reused bool
	}{
		{name: "max length accepted", header: strings.Repeat("a", 128), reused: true},
		{name: "over max length rejected", header: strings.Repeat("a", 129), reused: false},
		{name: "empty value not reused", header: "", reused: false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			router := newRequestIDRouter()

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/ping", nil)
			request.Header.Set(HeaderRequestID, testCase.header)
			router.ServeHTTP(recorder, request)

			got := recorder.Header().Get(HeaderRequestID)
			if !validRequestID.MatchString(got) {
				t.Fatalf("expected a valid request id in the response, got %q", got)
			}
			if testCase.reused && got != testCase.header {
				t.Fatalf("expected header %q to be reused, got %q", testCase.header, got)
			}
			if !testCase.reused && got == testCase.header {
				t.Fatalf("expected header %q to be replaced", testCase.header)
			}
		})
	}
}

func TestRequestIDExposedOnContext(t *testing.T) {
	router := newRequestIDRouter()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Body.String() != recorder.Header().Get(HeaderRequestID) {
		t.Fatalf("expected context id %q to match response header %q", recorder.Body.String(), recorder.Header().Get(HeaderRequestID))
	}
}

func TestRequestIDAppearsInAccessLog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logs bytes.Buffer
	router := gin.New()
	router.Use(RequestID(), gin.LoggerWithConfig(gin.LoggerConfig{Formatter: AccessLogFormatter, Output: &logs}))
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	request.Header.Set(HeaderRequestID, "client-123")
	router.ServeHTTP(recorder, request)

	if !strings.Contains(logs.String(), "request_id=client-123") {
		t.Fatalf("expected access log to contain request_id=client-123, got %q", logs.String())
	}
}

func newRequestIDRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, FromContext(ctx))
	})
	return router
}
