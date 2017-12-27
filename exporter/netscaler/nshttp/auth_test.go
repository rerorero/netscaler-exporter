package nshttp

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttpAuthSucceed(t *testing.T) {
	MockedNetscaler(HttpMockHandlers{}, func(ns *netscalerHttpImpl, server *httptest.Server) {
		err := ns.authorize()
		if err != nil {
			t.Fatalf("authorize() failed: %s", err)
		}

		uri, _ := url.Parse(ns.BaseHttpUrl() + "/v1")
		cookies := ns.http.Jar.Cookies(uri)
		if len(cookies) == 0 {
			t.Fatalf("No Set-Cookie")
		}
	})
}

func TestHttpAuthFailed(t *testing.T) {
	handlers := HttpMockHandlers{
		loginHandler: loginFailHandler,
	}
	MockedNetscaler(handlers, func(ns *netscalerHttpImpl, server *httptest.Server) {
		err := ns.authorize()
		if err == nil {
			t.Fatal("authorize() has not failed")
		}
	})
}
