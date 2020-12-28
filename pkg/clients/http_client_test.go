package clients

import (
	"encoding/json"
	"github.com/we4tech/uampnotif/pkg/testutils"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/we4tech/uampnotif/pkg/receivers"
)

const TestConfig = `
name: Test
id: test
request:
  valid_http_codes:
    - 201
    - 200
  params:
    - name: app_id
      label: App ID
    - name: api_key
      label: API Key
  method: POST
  url_tmpl: 'https://api.test.local/v2/applications/{{.FindParam "app_id"}}/deployments.json'
  headers:
    - name: Content-Type
      value_tmpl: application/json
    - name: X-Api-Key
      value_tmpl: '{{.FindParam "api_key"}}'
  body_tmpl: |
    { "deployment": { "revision": "{{.FindEnv "commit_hash"}}" } }
`

func TestNewRequest_ShouldCreateRequestWithoutError(t *testing.T) {
	params := map[string]string{
		"app_id":  "app-id-134",
		"api_key": "api-key-1345",
	}

	envVars := map[string]string{
		"commit_hash": "commit-hash-11333",
	}

	if _, err := NewHttpRequest(createIntSpec(), params, envVars, nil); err != nil {
		t.Errorf("could not send receivedRequest. Error - %s", err)
	}
}

func TestSendRequest_ShouldSuccessfullySendGetRequest(t *testing.T) {
	responseBody := []byte("hello world")
	mockClient := testutils.NewMockHttpClient(201, responseBody)

	request := buildCommonRequest(mockClient)
	resp, err := request.SendRequest()

	if err != nil {
		t.Errorf("could not complete receivedRequest. Error - %s", err)
	}

	if !resp.IsOK() {
		t.Errorf(
			"doesn't match http status code. Found: %d, Expected: %v",
			resp.Code, resp.validCodes)
	}

	if resp.Body != string(responseBody) {
		t.Errorf(
			"doesn't match sentResponse body. Found: %s, Expected: %s",
			resp.Body, responseBody)
	}
}

func TestSendRequest_ShouldSendHeaders(t *testing.T) {
	mockClient := testutils.NewMockHttpClient(201, []byte("hello"))

	request := buildCommonRequest(mockClient)
	resp, _ := request.SendRequest()

	if !resp.IsOK() {
		t.Error()
	}

	expectedHeaders := []string{"Content-Type", "X-Api-Kxx"}

	for key := range mockClient.ReceivedRequest.Header {
		if len(expectedHeaders) == sort.SearchStrings(expectedHeaders, key) {
			t.Errorf("could not find header - %s", key)
		}
	}
}

func TestSendRequest_ShouldSendPostBody(t *testing.T) {
	postBody := []byte("hello post body")
	mockClient := testutils.NewMockHttpClient(201, postBody)

	request := buildCommonRequest(mockClient)
	resp, _ := request.SendRequest()

	if !resp.IsOK() {
		t.Error()
	}

	bodyBytes, _ := ioutil.ReadAll(mockClient.ReceivedRequest.Body)
	bodyString := string(bodyBytes)

	if bodyString == "" {
		t.Error("could not find the matching post body")
	}

	bodyJson := &map[string]interface{}{}

	if err := json.Unmarshal(bodyBytes, bodyJson); err != nil {
		t.Errorf("could not parse body as a JSON. Err - %s", err)
	}

	deployment := (*bodyJson)["deployment"].(map[string]interface{})

	if deployment["revision"] != "commit-hash-123" {
		t.Error("could not find matching revision")
	}
}

func TestSendRequest_ShouldRaiseClientRequestErrorIfRequestFailed(t *testing.T) {
	postBody := []byte("hello post body")
	mockClient := testutils.
		NewMockHttpClient(500, postBody).
		RaiseError("unwanted error")

	request := buildCommonRequest(mockClient)
	_, err := request.SendRequest()

	if err == nil {
		t.Error("could not find the error")
	}

	if reflect.TypeOf(err) != reflect.TypeOf(clientRequestError{}) {
		t.Error("could not find the expected error class")
	}
}

func createIntSpec() *receivers.Spec {
	spec, _ := receivers.NewSpec([]byte(TestConfig))

	return spec
}

func buildCommonRequest(mockClient *testutils.MockHttpClient) Client {
	params := map[string]string{
		"app_id":  "app-id-123",
		"api_key": "api-key-1345",
	}
	envVars := map[string]string{
		"commit_hash": "commit-hash-123",
	}

	request, _ := NewHttpRequest(createIntSpec(), params, envVars, log.New(
		os.Stdout,
		"[test] ",
		log.Ltime,
	))
	request.SetClient(mockClient)

	return request
}
