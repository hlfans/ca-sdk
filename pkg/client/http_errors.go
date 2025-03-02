package client

import (
	"fmt"
	"strings"

	"github.com/hlfans/ca-sdk/pkg/response"
)

type ResponseError struct {
	Errors   []response.Message
	Messages []response.Message
}

func (err ResponseError) Error() string {
	return fmt.Sprintf("CA response error messages: %s", err.joinErrors())
}

func (err ResponseError) joinErrors() string {
	mes := make([]string, len(err.Errors))
	for i, m := range err.Errors {
		mes[i] = m.Message
	}

	return strings.Join(mes, `,`)
}

type ErrUnexpectedHTTPStatus struct {
	Status int
	Body   []byte
}

func (err ErrUnexpectedHTTPStatus) Error() string {
	return fmt.Sprintf("unexpected HTTP status code: %d with body %s", err.Status, string(err.Body))
}
