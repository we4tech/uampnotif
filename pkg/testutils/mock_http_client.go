package testutils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

type MockHttpClient struct {
	ReceivedRequest *http.Request

	responseCode int
	responseBody []byte

	raiseError string
}

func (mhc *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	mhc.ReceivedRequest = req

	if mhc.raiseError != "" {
		return nil, errors.New(mhc.raiseError)
	}

	body := ioutil.NopCloser(bytes.NewReader(mhc.responseBody))

	return &http.Response{
		StatusCode: mhc.responseCode,
		Body:       body,
	}, nil
}

func (mhc *MockHttpClient) RaiseError(err string) *MockHttpClient {
	mhc.raiseError = err

	return mhc
}

func NewMockHttpClient(rspCode int, rspBody []byte) *MockHttpClient {
	return &MockHttpClient{responseCode: rspCode, responseBody: rspBody}
}
