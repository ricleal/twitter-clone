package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ricleal/twitter-clone/internal/api/openapi"
	"github.com/rs/zerolog/log"
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
	json.NewEncoder(w).Encode(apiErr)
}
