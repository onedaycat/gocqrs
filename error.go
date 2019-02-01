package gocqrs

import (
	"github.com/onedaycat/errors"
)

var ErrNotFound = errors.NotFound("es1", "Not Found")
var ErrVersionInconsistency = errors.BadRequest("es2", "Version is inconsistency")
var ErrEncodingNotSupported = errors.InternalError("es3", "Unable unmarshal payload no support encoding")
