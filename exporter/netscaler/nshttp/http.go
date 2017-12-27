package nshttp

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/pkg/errors"
)

type NetscalerHttp interface {
	GetStats() (*NetscalerHttpStats, error)
}

type NetscalerHttpStats struct {
	VServers   map[string]*HttpVServerStats
	HttpHealth bool
}

type netscalerHttpImpl struct {
	http     *http.Client
	host     string
	httpPort int
	username string
	password string
}

func NewNetscalerHttpClient(host string, httpPort int, username string, password string, timeoutSec int) (NetscalerHttp, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize cookie")
	}

	httpClient := &http.Client{
		Jar:     jar,
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	return &netscalerHttpImpl{
		http:     httpClient,
		host:     host,
		httpPort: httpPort,
		username: username,
		password: password,
	}, nil
}

func (ns *netscalerHttpImpl) BaseHttpUrl() string {
	return fmt.Sprintf("http://%s:%d/nitro", ns.host, ns.httpPort)
}

func (ns *netscalerHttpImpl) GetStats() (*NetscalerHttpStats, error) {
	vservers, err := ns.getHttpVserverStats()
	if err != nil {
		return nil, err
	}
	return &NetscalerHttpStats{
		VServers:   vservers,
		HttpHealth: true,
	}, nil
}
