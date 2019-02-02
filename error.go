package gocqrs

import (
	"github.com/onedaycat/errors"
)

var ErrNotFound = errors.NotFound("es1", "Not Found")
var ErrVersionInconsistency = errors.BadRequest("es2", "Version is inconsistency")
var ErrEncodingNotSupported = errors.InternalError("es3", "Unable unmarshal payload, unsupport encoding")
var ErrEventLimitExceed = errors.InternalError("es4", "Number of events in aggregate limit exceed")
var ErrZeroVersionNotAllowed = errors.InternalError("es5", "Aggregate version should not be 0")
