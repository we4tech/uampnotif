package clients

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/we4tech/uampnotif/pkg/integrations"
	"github.com/we4tech/uampnotif/pkg/templates"
	"github.com/we4tech/uampnotif/pkg/validators"
)

const (
	GetMethod         = "GET"
	MaxRequestTimeout = 5 * time.Second
)

//
// clientValidationError represents an error that triggered while sending a
// receivedRequest to integration spec.
//
type clientValidationError struct {
	Errors validators.ValidationErrors
}

func (cve clientValidationError) Error() string {
	return fmt.Sprintf("ClientValidationErrors: %s", cve.Errors)
}

//
// clientRequestError represents an error triggered while sending client receivedRequest.
//
type clientRequestError struct {
	errorAt string
	err     error
}

func (cre clientRequestError) Error() string {
	return fmt.Sprintf(
		"ClientRequestError: Failed to process at %s. Error - %s",
		cre.errorAt, cre.err)
}

type httpClient struct {
	client  ClientImpl
	request *http.Request

	spec    *integrations.Spec
	tmplCtx *templates.TemplateContext

	url     string
	method  string
	headers map[string]string
	body    string
}

func (c *httpClient) SendRequest() (*Response, error) {
	if err := c.createRequest(); err != nil {
		return nil, err
	}

	c.addHeaders()
	c.addUserAgent()

	return c.execute()
}

func (c *httpClient) SetClient(clientImpl ClientImpl) {
	c.client = clientImpl
}

func (c *httpClient) getClient() ClientImpl {
	if c.client == nil {
		c.client = &http.Client{Timeout: MaxRequestTimeout}
	}

	return c.client
}

func (c *httpClient) validate() (bool, validators.ValidationErrors) {
	v := validators.NewValidator(c.spec, c.tmplCtx)

	if v.Validate() {
		return true, nil
	}

	return false, v.GetErrors()
}

//
// SendHttpRequest creates a *httpClient* with all required parameters to make
// a successful receivedRequest.
//
func NewHttpRequest(
	spec *integrations.Spec,
	parameters map[string]string,
	envVars map[string]string) (Client, error) {

	client := &httpClient{
		tmplCtx: templates.NewTemplateContext(parameters, envVars),
		spec:    spec,
	}

	if valid, errors := client.validate(); !valid {
		return nil, clientValidationError{Errors: errors}
	}

	if url, err := spec.Request.Url(client.tmplCtx); err != nil {
		return nil, clientRequestError{errorAt: "Url", err: err}
	} else {
		client.url = url
	}

	if err := client.buildHeaders(); err != nil {
		return nil, err
	}

	if body, err := spec.Request.Body(client.tmplCtx); err != nil {
		return nil, clientRequestError{errorAt: "Body", err: err}
	} else {
		client.body = body
	}

	return client, nil
}

func (c *httpClient) getPartialUrl() string {
	rx := regexp.MustCompile("(.+)[?/].+")

	return rx.FindStringSubmatch(c.url)[0]
}

func (c *httpClient) getBodyAsReader() io.Reader {
	if strings.ToUpper(c.method) == GetMethod ||
		c.body == "" {
		return nil
	}

	return bytes.NewReader([]byte(c.body))
}

func (c *httpClient) addHeaders() {
	if len(c.headers) > 0 {
		for key, value := range c.headers {
			c.request.Header.Add(key, value)
		}
	}
}

func (c *httpClient) addUserAgent() {
	c.request.Header.Add("User-Agent", "uampnotif")
}

func (c *httpClient) execute() (*Response, error) {
	if res, err := c.getClient().Do(c.request); err != nil {
		return nil, clientRequestError{errorAt: "Request.Received", err: err}
	} else {
		if body, err := ioutil.ReadAll(res.Body); err != nil {
			return nil, clientRequestError{errorAt: "Request.ReadingBody", err: err}
		} else {
			log.Printf(
				"sendGetRequest: Received sentResponse for method:%s from "+
					"url:%s status:%d",
				c.method, c.getPartialUrl(), res.StatusCode)

			return &Response{
				Code:       res.StatusCode,
				Body:       string(body),
				validCodes: c.spec.Request.ValidHttpCodes,
			}, nil
		}
	}
}

func (c *httpClient) createRequest() error {
	log.Printf(
		"sendPostRequest: Sending method:%s to url:%s",
		c.method, c.getPartialUrl())

	request, err := http.NewRequest(
		strings.ToUpper(c.method), c.url, c.getBodyAsReader())

	if err != nil {
		return clientRequestError{errorAt: "Request.Send", err: err}
	}

	c.request = request

	return nil
}

func (c *httpClient) buildHeaders() error {
	c.headers = make(map[string]string)

	err := c.spec.Request.Headers.ForEach(
		func(i int, pHeader integrations.ParsedHeader) (bool, error) {
			if value, err := pHeader.GetValue(c.tmplCtx); err != nil {
				return true, clientRequestError{
					errorAt: fmt.Sprintf("Request.Headers.%s", pHeader.GetName()),
					err:     err}
			} else {
				c.headers[pHeader.GetName()] = value
			}

			return false, nil
		})

	if err != nil {
		return clientRequestError{errorAt: "Headers", err: err}
	}

	return nil
}
