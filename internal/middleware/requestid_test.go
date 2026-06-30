package middleware

import (
	"net/http"
	"net/http/httptest"
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

func TestRequestIDExposedOnContext(t *testing.T) {
	router := newRequestIDRouter()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Body.String() != recorder.Header().Get(HeaderRequestID) {
		t.Fatalf("expected context id %q to match response header %q", recorder.Body.String(), recorder.Header().Get(HeaderRequestID))
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
