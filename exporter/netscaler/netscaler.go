package netscaler

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/pkg/errors"
)

type Netscaler interface {
	Authorize() error
	GetHttpVserverStats() ([]HttpVServerStats, error)
	BaseHttpUrl() string
}

type netscalerImpl struct {
	http     *http.Client
	host     string
	httpPort int
	username string
	password string
}

func NewNetscalerClient(host string, httpPort int, username string, password string) (Netscaler, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize cookie")
	}

	httpClient := &http.Client{
		Jar: jar,
	}

	return &netscalerImpl{
		http:     httpClient,
		host:     host,
		httpPort: httpPort,
		username: username,
		password: password,
	}, nil
}

func (ns *netscalerImpl) BaseHttpUrl() string {
	return fmt.Sprintf("http://%s:%d/nitro", ns.host, ns.httpPort)
}
