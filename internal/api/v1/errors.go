package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ricleal/twitter-clone/internal/api/v1/openapi"
)

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendAPIError(ctx context.Context, w http.ResponseWriter, code int, message string, err error) {
	log.Ctx(ctx).Error().Err(err).Msg(message)
	apiErr := openapi.Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiErr) //nolint:errcheck //ignore error
}
