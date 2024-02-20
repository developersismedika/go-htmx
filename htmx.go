// Package htmx offers a streamlined integration with HTMX in Go applications.
// It implements the standard io.Writer interface and includes middleware support, but it is not required.
// Allowing for the effortless incorporation of HTMX features into existing Go applications.
package htmx

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	DefaultSwapDuration = time.Duration(0 * time.Millisecond)
	DefaultSettleDelay  = time.Duration(20 * time.Millisecond)

	DefaultNotificationKey = "showMessage"
)

type (
	HTMX struct {
		log *zap.Logger
	}
)

// New returns a new htmx instance.
func New() *HTMX {
	// prepares a basic logger
	log, _ := zap.NewProduction(zap.WithCaller(false))

	return &HTMX{
		log: log,
	}
}

// SetLog sets the logger for the htmx instance.
func (h *HTMX) SetLog(log *zap.Logger) {
	h.log = log
}

// NewHandler returns a new htmx handler.
func (h *HTMX) NewHandler(w http.ResponseWriter, r *http.Request) *Handler {
	return &Handler{
		w:        w,
		r:        r,
		request:  h.HxHeader(r),
		response: h.HxResponseHeader(w.Header()),
		log:      h.log,
	}
}

// IsHxRequest returns true if the request is a htmx request.
func IsHxRequest(r *http.Request) bool {
	return HxStrToBool(r.Header.Get(HxRequestHeaderRequest.String()))
}

// IsHxBoosted returns true if the request is a htmx request and the request is boosted
func IsHxBoosted(r *http.Request) bool {
	return HxStrToBool(r.Header.Get(HxRequestHeaderBoosted.String()))
}

// IsHxHistoryRestoreRequest returns true if the request is a htmx request and the request is a history restore request
func IsHxHistoryRestoreRequest(r *http.Request) bool {
	return HxStrToBool(r.Header.Get(HxRequestHeaderHistoryRestoreRequest.String()))
}

// RenderPartial returns true if the request is an HTMX request that is either boosted or a hx request,
// provided it is not a history restore request.
func RenderPartial(r *http.Request) bool {
	return (IsHxRequest(r) || IsHxBoosted(r)) && !IsHxHistoryRestoreRequest(r)
}

// HxStrToBool converts a string to a boolean value.
func HxStrToBool(str string) bool {
	return strings.EqualFold(str, "true")
}

// HxBoolToStr converts a boolean value to a string.
func HxBoolToStr(b bool) string {
	if b {
		return "true"
	}

	return "false"
}
