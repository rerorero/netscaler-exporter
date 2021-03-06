package nshttp

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestGetVserverStatsSucceed(t *testing.T) {
	MockedNetscaler(HttpMockHandlers{}, func(ns *netscalerHttpImpl, server *httptest.Server) {
		expected := map[string]*json.RawMessage{}
		err := json.Unmarshal([]byte(statsBody), &expected)
		if err != nil {
			t.Fatalf("Invalid json")
		}

		stats, err := ns.getHttpVserverStats()
		if err != nil {
			t.Fatalf("GetVserverStats() failed: %s\n", err)
		}

		if len(stats) == len(expected) {
			t.Fatalf("Length is not matched %v - %v\n", stats, expected)
		}
	})
}

func TestGetVserverStatsFailedInAuth(t *testing.T) {
	handlers := HttpMockHandlers{
		loginHandler: loginFailHandler,
	}
	MockedNetscaler(handlers, func(ns *netscalerHttpImpl, server *httptest.Server) {
		_, err := ns.getHttpVserverStats()
		if err == nil {
			t.Fatal("GetHttpVserverStats() succeed")
		}
	})
}
