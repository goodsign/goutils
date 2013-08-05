package jsonrpc

import (
	"bytes"
	"github.com/goodsign/rpc/json"
	"io/ioutil"
	"net/http"
	"time"
)

// ServiceClient contains data needed to connect to gorilla-powered json-rpc services
type ServiceClient struct {
	address          string
	callRetryCount   int
	callRetryTimeout time.Duration
	client           *http.Client
}

// NewServiceClient creates a new service client to work with the
// service with a specified address (full address with port).
func NewServiceClient(addr string, retryCount int, retryTimeout time.Duration) (*ServiceClient, error) {
	c := new(ServiceClient)
	c.address = addr
	c.callRetryTimeout = retryTimeout
	c.callRetryCount = retryCount
	c.client = http.DefaultClient
	return c, nil
}

func (c *ServiceClient) doRetry(method string, bodyBytes []byte, address string, res interface{}) (e error) {
	rtc := c.callRetryCount

	for {
		var req *http.Request
		req, e = http.NewRequest("POST", address, bytes.NewReader(bodyBytes))

		if e == nil {
			req.Header.Set("Content-Type", "application/json")
			var resp *http.Response
			resp, e = c.client.Do(req)

			if e == nil {
				var bodyb []byte
				bodyb, e = ioutil.ReadAll(resp.Body)

				if e == nil {
					serr := json.DecodeClientResponse(bytes.NewBuffer(bodyb), res)

					if serr != nil {
						logger.Errorf("Error while decoding '%s' :'%s'", string(bodyb), serr)
					}
					return serr
				} else {
					logger.Criticalf("E2: %s", e)
				}
			}
		}

		if e == nil || rtc <= 0 {
			logger.Errorf("Leaving after all retries with error: '%s'", e)
			return
		}

		logger.Warn(e)
		logger.Tracef("Reconnect (#%d left) for '%s'", rtc, method)
		rtc--
		time.Sleep(c.callRetryTimeout)
	}

	return
}

// GetResult returns a result of a method call to a json-rpc service with specified arguments.
// It performs reconnects according to retry args passed in NewServiceClient call.
func (c *ServiceClient) GetResult(method string, args, res interface{}) error {
	logger.Debugf("JSON-RPC Client call: '%s'", method)

	// Convert client request and params to correct encoded bytes for gorilla
	// according to JSON-RPC 1.0: http://json-rpc.org/wiki/specification
	buf, err := json.EncodeClientRequest(method, args)

	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Tracef("Post request to '%s'", c.address)

	// Get result and decode it
	err = c.doRetry(method, buf, c.address, res)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
