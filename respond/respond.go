package respond

import (
	"net/http"

	"github.com/sermodigital/json"
)

// Common MIME types.
const (
	HTMLMIME = "text/html; charset=utf-8"
	JSONMIME = "application/json; charset=utf-8"
	TextMIME = "text/plain; charset=utf-8"

	ContentType   = "Content-Type"
	AccessControl = "Access-Control-Allow-"
)

const APIVersion = "0.1"

// Response is a response from the SermoCRM API.
type Response struct {
	// Arbitrary response data.
	Data interface{} `json:"data,omitempty"`

	// Non-empty if an error occurred.
	Error string `json:"error,omitempty"`
}

// Write writes the Response to w. It sends the Response in a similar format,
// except it adds the API version as well as a unique identifier for the
// specific API response.
func (r Response) Write(w http.ResponseWriter) {
	err := json.MarshalStream(w, r)
	if err != nil {
		SystemError(w, nil)
	}
}

// http://stackoverflow.com/a/2669766/2967113
// https://docs.angularjs.org/api/ng/service/$http#json-vulnerability-protection
var dontBeEvil = []byte(")]}',\n")

func (r Response) json(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	w.Header().Set(ContentType, JSONMIME)

	// Don't be evil.
	w.Write(dontBeEvil)
	r.Write(w)
}

func OK(w http.ResponseWriter, data interface{}) {
	(Response{Data: data}).json(w, http.StatusOK)
}

// BadMethod reponds with http.StatusMethodNotAllowed.
func BadMethod(w http.ResponseWriter, _ error) {
	(Response{Error: "method not allowed"}.json(w, http.StatusMethodNotAllowed))
}

// BadRequest sends an error response to with http.StatusNotAcceptable.
func BadRequest(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, http.StatusBadRequest)
}

// Forbidden sends an error response with http.StatusForbidden.
func Forbidden(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, http.StatusForbidden)
}

// Unauthenticated sends an error response with http.StatusUnauthorized.
func Unauthenticated(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, http.StatusUnauthorized)
}

// Unavailable sends an error response with http.StatusServiceUnavailable.
func Unavailable(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, http.StatusServiceUnavailable)
}

// OutOfRange sends an error response with http.StatusRequestedRangeNotSatisfiable.
func OutOfRange(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(
		w, http.StatusRequestedRangeNotSatisfiable)
}

// Exhausted sends an error response with http.StatusTooManyRequests.
func Exhausted(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, http.StatusTooManyRequests)
}

// Unimplemented sends an error response with http.StatusTooManyRequests.
func Unimplemented(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, http.StatusNotImplemented)
}

// InternalServerError sends an error response with http.StatusInternalServerError.
func InternalServerError(w http.ResponseWriter, err error) {
	(Response{Error: "internal server error"}).json(w, http.StatusInternalServerError)
}

// NotFound sends an error response with http.StatusNotFound.
func NotFound(w http.ResponseWriter, _ error) {
	(Response{Error: "resource not found"}).json(w, http.StatusNotFound)
}

const (
	StatusUnknown            = http.StatusInternalServerError
	StatusCanceled           = http.StatusRequestTimeout
	StatusDeadlineExceeded   = http.StatusRequestTimeout
	StatusAlreadyExists      = http.StatusConflict
	StatusFailedPrecondition = http.StatusPreconditionFailed
	StatusAborted            = http.StatusConflict
	StatusDataLoss           = http.StatusInternalServerError
)

// Unknown sends an error response with StatusUnknown.
func Unknown(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, StatusUnknown)
}

// Canceled sends an error response with StatusCanceled.
func Canceled(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, StatusCanceled)
}

// DeadlineExceeded sends an error response with StatusDeadlineExceeded.
func DeadlineExceeded(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, StatusDeadlineExceeded)
}

// AlreadyExists sends an error response with StatusAlreadyExists.
func AlreadyExists(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, StatusAlreadyExists)
}

// FailedPrecondition sends an error response with StatusFailedPrecondition.
func FailedPrecondition(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, StatusFailedPrecondition)
}

// Aborted sends an error response with StatusAborted.
func Aborted(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, StatusAborted)
}

// Unknown sends an error response with StatusDataLoss.
func DataLoss(w http.ResponseWriter, err error) {
	(Response{Error: err.Error()}).json(w, StatusDataLoss)
}

// SystemError responds with a generic error response if marshaling an
// APIResponse fails.
func SystemError(w http.ResponseWriter, _ error) {
	const genericError = `{"error": "SYSTEM ERROR"}`
	w.Header().Set(ContentType, JSONMIME)
	w.Write([]byte(genericError))
}
