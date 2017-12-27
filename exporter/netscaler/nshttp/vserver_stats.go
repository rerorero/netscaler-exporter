package nshttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type HttpVServerStats struct {
	Name                  string  `json:"name"`
	VserverSurgeCount     float64 `json:"vsvrsurgecount,string"`
	EstablishedConn       float64 `json:"establishedconn,string"`
	InactiveServices      float64 `json:"inactsvcs,string"`
	Health                float64 `json:"vslbhealth,string"`
	PrimaryIpAddress      string  `json:"primaryipaddress"`
	PrimaryPort           int     `json:"primaryport"`
	Type                  string  `json:"type"`
	State                 string  `json:"state"`
	ActiveServices        float64 `json:"actsvcs,string"`
	TotalHits             float64 `json:"tothits,string"`
	HitsRate              float64 `json:"hitsrate"`
	TotalRequests         float64 `json:"totalrequests,string"`
	RequestsRate          float64 `json:"requestsrate"`
	TotalResponses        float64 `json:"totalresponses,string"`
	ResponsesRate         float64 `json:"responsesrate"`
	TotalRequestBytes     float64 `json:"totalrequestbytes,string"`
	RequestBytesRate      float64 `json:"requestbytesrate"`
	TotalResponseBytes    float64 `json:"totalresponsebytes,string"`
	ResponseBytesRate     float64 `json:"responsebytesrate"`
	TotalPackateReceived  float64 `json:"totalpktsrecvd,string"`
	PackageReceivedRate   float64 `json:"pktsrecvdrate"`
	TotalPackageSent      float64 `json:"totalpktssent,string"`
	PackageSentRate       float64 `json:"pktssentrate"`
	SurgeCount            float64 `json:"surgecount,string"`
	ServiceSurgeCount     float64 `json:"svcsurgecount,string"`
	InvlidRequestResponse float64 `json:"invalidrequestresponse,string"`
}
type NsVserversStatResult struct {
	Stats []HttpVServerStats `json:"lbvserver"`
}

// ref. https://docs.citrix.com/en-us/netscaler/11/nitro-api/nitro-rest/nitro-rest-general/nitro-rest-statistics.html
func (ns *netscalerHttpImpl) getHttpVserverStats() (map[string]*HttpVServerStats, error) {
	req, err := http.NewRequest(http.MethodGet, ns.BaseHttpUrl()+"/v1/stat/lbvserver", nil)
	if err != nil {
		return nil, errors.Wrap(err, "New request failed")
	}
	result := &NsVserversStatResult{}

	err = ns.withAuth(withAuthParam{
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

	m := map[string]*HttpVServerStats{}
	for _, stat := range result.Stats {
		m[stat.Name] = &stat
	}

	return m, nil
}
