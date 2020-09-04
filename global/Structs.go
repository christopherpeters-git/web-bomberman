package global

//Type of error that has enough information to answer a http request
type DetailedHttpError struct {
	status      int    //HTTP status code
	publicError string //Error message that may be exposed to the client
	error       string //Error message that should only appear in the logs
}

func NewDetailedHttpError(status int, publicError string, error string) *DetailedHttpError {
	return &DetailedHttpError{
		status:      status,
		publicError: publicError,
		error:       error,
	}
}

func (r *DetailedHttpError) PublicError() string {
	return r.publicError
}
func (r *DetailedHttpError) Error() string {
	return r.error
}
func (r *DetailedHttpError) Status() int {
	return r.status
}
