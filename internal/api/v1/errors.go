package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ricleal/twitter-clone/internal/api/v1/openapi"
)

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendAPIError(ctx context.Context, w http.ResponseWriter, code int, message string, err error) {
	slog.ErrorContext(ctx, message, "error", err)
	apiErr := openapi.Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiErr) //nolint:errcheck //ignore error
}
