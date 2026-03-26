package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ricleal/twitter-clone/internal/api/v1/openapi"
)

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendAPIError(
	logger *slog.Logger,
	w http.ResponseWriter,
	code int,
	message string,
	err error,
) {
	logger.Error(message, "error", err)
	apiErr := openapi.Error{
		Code:    int32(code), //nolint:gosec // HTTP status codes safely fit in int32
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiErr) //nolint:errcheck,gosec //ignore error
}
