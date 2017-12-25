package netscaler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type VServerStats struct {
	Name string `json:name`
}
type NsVserversStatResult struct {
	Stats []VServerStats `json:lbvserver`
}

// ref. https://docs.citrix.com/en-us/netscaler/11/nitro-api/nitro-rest/nitro-rest-general/nitro-rest-statistics.html
func (ns *netscalerImpl) GetVserverStats() ([]VServerStats, error) {
	req, err := http.NewRequest("GET", ns.baseHttpUrl()+"/v1/stat/lbvserver", nil)
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
