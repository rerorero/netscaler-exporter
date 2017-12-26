package netscaler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type HttpVServerStats struct {
	Name                  string  `json:"name"`
	VserverSrugeCount     int     `json:"vsvrsurgecount,string"`
	EstablishedConn       int     `json:"establishedconn,string"`
	InactiveServices      int     `json:"inactsvcs,string"`
	Health                float64 `json:"vslbhealth,string"`
	PrimaryIpAddress      string  `json:"primaryipaddress"`
	PrimaryPort           int     `json:"primaryport"`
	Type                  string  `json:"type"`
	State                 string  `json:"state"`
	ActiveServices        int     `json:"actsvcs,string"`
	TotalHits             int64   `json:"tothits,string"`
	HitsRate              int     `json:"hitsrate"`
	TotalRequests         int64   `json:"totalrequests,string"`
	RequestsRate          int     `json:"requestsrate"`
	TotalResponses        int64   `json:"totalresponses,string"`
	ResponsesRate         int     `json:"responsesrate"`
	TotalRequestBytes     int64   `json:"totalrequestbytes,string"`
	RequestBytesRate      int     `json:"requestbytesrate"`
	TotalResponseBytes    int64   `json:"totalresponsebytes,string"`
	ResponseBytesRate     int     `json:"responsebytesrate"`
	TotalPackateReceived  int64   `json:"totalpktsrecvd,string"`
	PackageReceivedRate   int     `json:"pktsrecvdrate"`
	TotalPackageSent      int64   `json:"totalpktssent,string"`
	PackageSentRate       int     `json:"pktssentrate"`
	SurgeCount            int     `json:"surgecount,string"`
	ServiceSurgeCount     int     `json:"svcsurgecount,string"`
	InvlidRequestResponse int64   `json:"invalidrequestresponse,string"`
}
type NsVserversStatResult struct {
	Stats []HttpVServerStats `json:"lbvserver"`
}

// ref. https://docs.citrix.com/en-us/netscaler/11/nitro-api/nitro-rest/nitro-rest-general/nitro-rest-statistics.html
func (ns *netscalerImpl) GetHttpVserverStats() ([]HttpVServerStats, error) {
	req, err := http.NewRequest(http.MethodGet, ns.BaseHttpUrl()+"/v1/stat/lbvserver", nil)
	if err != nil {
		return nil, errors.Wrap(err, "New request failed")
	}
	result := &NsVserversStatResult{}

	err = ns.WithAuth(withAuthParam{
		req: req,
		f: func(resp *http.Response, body []byte) error {
			err2 := json.Unmarshal(body, result)
			if err2 != nil {
				return errors.Wrap(err2, fmt.Sprintf("Failed to parse vserver stats response body, %s", string(body)))
			}
			return nil
		},
	})

	if err != nil {
		return nil, err
	}
	return result.Stats, nil
}
