package netscaler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
)

type HttpMock interface {
	Mocked(func(Netscaler, *httptest.Server))
}

type HttpMockHandlers struct {
	loginHandler http.HandlerFunc
	statsHandler http.HandlerFunc
	rootHandler  http.HandlerFunc
}

type httpMockImpl struct {
	handlers HttpMockHandlers
}

var (
	loginSuccessHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		cookie := &http.Cookie{
			Name:  "NITRO_AUTH_TOKEN",
			Value: "DEADBEEF",
			Path:  "/nitro/v1",
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "")
	})

	loginFailHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "")
	})

	notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	statsBody = `
			{
			"errorcode": 0,
			"message": "Done",
			"severity": "NONE",
			"lbvserver": [
				{
				"name": "lbvserver1",
				"sortorder": "descending",
				"vsvrsurgecount": "0",
				"establishedconn": "0",
				"inactsvcs": "3",
				"vslbhealth": "0",
				"primaryipaddress": "192.168.0.10",
				"primaryport": 25,
				"type": "TCP",
				"state": "DOWN",
				"actsvcs": "0",
				"tothits": "0",
				"hitsrate": 0,
				"totalrequests": "0",
				"requestsrate": 0,
				"totalresponses": "0",
				"responsesrate": 0,
				"totalrequestbytes": "0",
				"requestbytesrate": 0,
				"totalresponsebytes": "0",
				"responsebytesrate": 0,
				"totalpktsrecvd": "0",
				"pktsrecvdrate": 0,
				"totalpktssent": "0",
				"pktssentrate": 0,
				"curclntconnections": "0",
				"cursrvrconnections": "0",
				"surgecount": "0",
				"svcsurgecount": "0",
				"sothreshold": "0",
				"totspillovers": "0",
				"labelledconn": "0",
				"pushlabel": "0",
				"deferredreq": "0",
				"deferredreqrate": 0,
				"invalidrequestresponse": "0",
				"invalidrequestresponsedropped": "0",
				"totvserverdownbackuphits": "0",
				"curmptcpsessions": "0",
				"cursubflowconn": "0"
				},
				{
				"name": "lbvserver2",
				"sortorder": "descending",
				"vsvrsurgecount": "0",
				"establishedconn": "0",
				"inactsvcs": "0",
				"vslbhealth": "100",
				"primaryipaddress": "192.168.0.20",
				"primaryport": 80,
				"type": "HTTP",
				"state": "UP",
				"actsvcs": "1",
				"tothits": "141205",
				"hitsrate": 0,
				"totalrequests": "141205",
				"requestsrate": 0,
				"totalresponses": "141205",
				"responsesrate": 0,
				"totalrequestbytes": "49944206",
				"requestbytesrate": 0,
				"totalresponsebytes": "44559139",
				"responsebytesrate": 0,
				"totalpktsrecvd": "141209",
				"pktsrecvdrate": 0,
				"totalpktssent": "141215",
				"pktssentrate": 0,
				"curclntconnections": "0",
				"cursrvrconnections": "0",
				"surgecount": "0",
				"svcsurgecount": "0",
				"sothreshold": "0",
				"totspillovers": "0",
				"labelledconn": "0",
				"pushlabel": "0",
				"deferredreq": "0",
				"deferredreqrate": 0,
				"invalidrequestresponse": "0",
				"invalidrequestresponsedropped": "0",
				"totvserverdownbackuphits": "0",
				"curmptcpsessions": "0",
				"cursubflowconn": "0"
				},
			{
				"name": "lbvserver3",
				"sortorder": "descending",
				"vsvrsurgecount": "0",
				"establishedconn": "55",
				"inactsvcs": "0",
				"vslbhealth": "100",
				"primaryipaddress": "192.168.0.30",
				"primaryport": 22,
				"type": "TCP",
				"state": "UP",
				"actsvcs": "1",
				"tothits": "55876398",
				"hitsrate": 2,
				"totalrequests": "0",
				"requestsrate": 0,
				"totalresponses": "0",
				"responsesrate": 0,
				"totalrequestbytes": "9261622292",
				"requestbytesrate": 442,
				"totalresponsebytes": "1256464744519",
				"responsebytesrate": 90081,
				"totalpktsrecvd": "2533119072",
				"pktsrecvdrate": 54,
				"totalpktssent": "2759605964",
				"pktssentrate": 79,
				"curclntconnections": "55",
				"cursrvrconnections": "55",
				"surgecount": "0",
				"svcsurgecount": "0",
				"sothreshold": "0",
				"totspillovers": "0",
				"labelledconn": "0",
				"pushlabel": "0",
				"deferredreq": "0",
				"deferredreqrate": 0,
				"invalidrequestresponse": "0",
				"invalidrequestresponsedropped": "0",
				"totvserverdownbackuphits": "0",
				"curmptcpsessions": "0",
				"cursubflowconn": "0"
				}
			]}`

	statsDefaultHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.com.citrix.netscaler.lbvserver+json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, statsBody)
	})
)

func NewHttpMock(handlers HttpMockHandlers) HttpMock {
	if handlers.loginHandler == nil {
		handlers.loginHandler = loginSuccessHandler
	}

	if handlers.statsHandler == nil {
		handlers.statsHandler = statsDefaultHandler
	}

	if handlers.rootHandler == nil {
		handlers.rootHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/nitro/v1/config/login" && r.Method == http.MethodPost {
				handlers.loginHandler(w, r)

			} else if r.URL.Path == "/nitro/v1/stat/lbvserver" && r.Method == http.MethodGet {
				if verifyCookie(r) {
					handlers.statsHandler(w, r)
				} else {
					loginFailHandler(w, r)
				}
			} else {
				notFoundHandler(w, r)
			}
		})
	}
	return &httpMockImpl{
		handlers: handlers,
	}
}

func verifyCookie(r *http.Request) bool {
	cookie, err := r.Cookie("NITRO_AUTH_TOKEN")
	return err == nil && cookie != nil
}

func (mock *httpMockImpl) Mocked(f func(Netscaler, *httptest.Server)) {
	server := httptest.NewServer(mock.handlers.rootHandler)
	defer server.Close()

	uri, _ := url.Parse(server.URL)
	port, _ := strconv.Atoi(uri.Port())
	ns, err := NewNetscalerClient(uri.Hostname(), port, "Alice", "secret")
	if err != nil {
		log.Fatalf("Could not start mock server: %s", err)
	}
	f(ns, server)
}

func MockedNetscaler(handlers HttpMockHandlers, f func(Netscaler, *httptest.Server)) {
	mock := NewHttpMock(handlers)
	mock.Mocked(f)
}
