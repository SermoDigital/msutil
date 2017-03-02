package respond

import (
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sermodigital/errors"
	"github.com/sermodigital/json"
	"github.com/sermodigital/pools"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ErrFunc is a function that writes the given error to w.
type ErrFunc func(http.ResponseWriter, error)

var codeJumpTable = [...]ErrFunc{
	codes.OK:                 InternalServerError,
	codes.Canceled:           Canceled,
	codes.Unknown:            Unknown,
	codes.InvalidArgument:    BadRequest,
	codes.DeadlineExceeded:   DeadlineExceeded,
	codes.NotFound:           NotFound,
	codes.AlreadyExists:      AlreadyExists,
	codes.PermissionDenied:   Forbidden,
	codes.Unauthenticated:    Unauthenticated,
	codes.ResourceExhausted:  Exhausted,
	codes.FailedPrecondition: FailedPrecondition,
	codes.Aborted:            Aborted,
	codes.OutOfRange:         OutOfRange,
	codes.Unimplemented:      Unimplemented,
	codes.Internal:           InternalServerError,
	codes.Unavailable:        Unavailable,
	codes.DataLoss:           DataLoss,
}

func ErrorHandler(_ context.Context, _ runtime.Marshaler,
	w http.ResponseWriter, _ *http.Request, err error) {
	Any(w, err)
}

// Any determines the type of response to write to the http.ResponseWriter
// from the error's RPC code. err must != nil. Unknown errors are considered
// to be Internal Server Errors.
func Any(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	code := grpc.Code(err)
	if code < codes.Code(len(codeJumpTable)) {
		codeJumpTable[code](w, err)
	} else {
		InternalServerError(w,
			errors.Errorf("invalid gRPC error code: %#v", err))
	}
}

// T implements grpc-ecosystem/grpc-gateway/runtime.Marshaler.
type T struct{}

// ContentType implements grpc-ecosystem/grpc-gateway/runtime.Marshaler.
func (t T) ContentType() string {
	return "application/json"
}

// Marshal implements grpc-ecosystem/grpc-gateway/runtime.Marshaler.
func (t T) Marshal(v interface{}) ([]byte, error) {
	// empty.Empty means we don't have a return body.
	_, ok := v.(*empty.Empty)
	if ok {
		return nil, nil
	}

	b := pools.GetBuffer()
	_, err := b.Write(dontBeEvil)
	if err != nil {
		return nil, errors.Internal(err)
	}

	// For some odd reason grpc-gateway wraps the stream result in an object
	// with the key "result".
	pbm, ok := v.(map[string]proto.Message)
	if ok {
		v = pbm["result"]
	}
	err = json.MarshalStream(b, Response{Data: v})
	if err != nil {
		return nil, errors.Internal(err)
	}
	return b.UnsafeBytes(), nil
}

// Unmarshal implements grpc-ecosystem/grpc-gateway/runtime.Marshaler.
func (t T) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return errors.Internal(errors.New("v does not implement proto.Message"))
	}
	err := proto.Unmarshal(data, pb)
	if err != nil {
		if err == json.ErrTooLarge {
			return errors.InvalidArg(err)
		}
		return errors.Internal(err)
	}
	return nil
}

// NewDecoder implements grpc-ecosystem/grpc-gateway/runtime.Marshaler.
func (t T) NewDecoder(r io.Reader) runtime.Decoder {
	return json.NewDecoder(r)
}

// NewEncoder implements grpc-ecosystem/grpc-gateway/runtime.Marshaler.
func (t T) NewEncoder(w io.Writer) runtime.Encoder {
	return json.NewEncoder(w)
}

var _ runtime.Marshaler = T{}
