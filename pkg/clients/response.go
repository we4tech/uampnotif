package clients

//
// Response represents the http code and body content.
//
type Response struct {
	Code int
	Body string

	validCodes []int
}

//
// IsOk returns true if status code matches one of the validCodes
//
func (r *Response) IsOK() bool {
	for _, vCode := range r.validCodes {
		if vCode == r.Code {
			return true
		}
	}

	return false
}
