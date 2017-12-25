package netscaler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type withAuthParam struct {
	req       *http.Request
	f         func(*http.Response, []byte) error
	retryAuth bool
}

// ref. https://docs.citrix.com/en-us/netscaler/11/nitro-api/nitro-rest/nitro-rest-connecting.html
func (ns *netscalerImpl) Authorize() error {
	path := ns.baseHttpUrl() + "/v1/config/login"
	data := fmt.Sprintf(`{ 
		"login": { 
			"username":"%s", 
			"password\":"%s" 
		}
	}`, ns.username, ns.password)

	r, err := ns.http.Post(path, "application/vnd.com.citrix.netscaler.login+json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to connect to %s", ns.host))
	}

	if r != nil {
		defer r.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to login to %s", ns.host))
	}

	if (r.StatusCode % 100) != 2 {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to read response body, status=%d", r.StatusCode))
		}
		return fmt.Errorf("Failed to login to %s, status=%d, body=%s", ns.host, r.StatusCode, string(body))
	}

	return nil
}

func (ns *netscalerImpl) WithAuth(param withAuthParam) error {
	resp, err := ns.http.Do(param.req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to request to method=%s, url=%s", param.req.Method, param.req.URL))
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to request to method=%s, url=%s", param.req.Method, param.req.URL))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to read response body, method=%s, url=%s", param.req.Method, param.req.URL))
	}

	if (resp.StatusCode % 100) == 2 {
		return param.f(resp, body)
	} else if resp.StatusCode == 401 {
		// retry authorization once
		if param.retryAuth {
			err = ns.Authorize()
			if err != nil {
				return err
			}
			noAuth := param
			noAuth.retryAuth = false
			return ns.WithAuth(noAuth)
		}
		return fmt.Errorf("Failed to login to %s, status=%d, body=%s", ns.host, resp.StatusCode, string(body))
	} else {
		return fmt.Errorf("Respond error, method=%s, url=%s, status=%d, body=%s", param.req.Method, param.req.URL, resp.StatusCode, string(body))
	}
}
