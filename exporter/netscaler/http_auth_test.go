package netscaler

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttpAuthSucceed(t *testing.T) {
	MockedNetscaler(HttpMockHandlers{}, func(ns Netscaler, server *httptest.Server) {
		err := ns.Authorize()
		if err != nil {
			t.Fatalf("Authorize() failed: %s", err)
		}

		nsImpl := ns.(*netscalerImpl)
		uri, _ := url.Parse(ns.BaseHttpUrl() + "/v1")
		cookies := nsImpl.http.Jar.Cookies(uri)
		if len(cookies) == 0 {
			t.Fatalf("No Set-Cookie")
		}
	})
}

func TestHttpAuthFailed(t *testing.T) {
	handlers := HttpMockHandlers{
		loginHandler: loginFailHandler,
	}
	MockedNetscaler(handlers, func(ns Netscaler, server *httptest.Server) {
		err := ns.Authorize()
		if err == nil {
			t.Fatal("Authorize() has not failed")
		}
	})
}