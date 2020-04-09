package errors

// RestError represents a Rest HTTP Error that can be returned from a controller
type RestError struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	OriginalError error  `json:"-"`
}

func (re *RestError) Error() string {
	return re.Message
}

func (re *RestError) Unwrap() error {
	return re.OriginalError
}
