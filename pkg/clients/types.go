package clients

import "net/http"

type Client interface {
	//
	// SendRequest dispatches a receivedRequest based on the integration spec.
	//
	SendRequest() (*Response, error)

	//
	// SetClient sets the internal http client.
	//
	SetClient(client ClientImpl)
}

//
// ClientImpl provides an interface for a specific implementation.
//
// Added for convenience to create a mock implementation.
//
type ClientImpl interface {
	Do(req *http.Request) (*http.Response, error)
}
