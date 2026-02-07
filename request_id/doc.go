// Package request_id provides context-based request ID propagation.
//
// It enables distributed tracing and request correlation across service boundaries.
//
// Example:
//
//	// Set in middleware
//	ctx = request_id.ContextWithRequestID(r.Context(), r.Header.Get("X-Request-ID"))
//
//	// Retrieve anywhere
//	requestID := request_id.RequestIDFromContext(ctx)
//	log.Info(ctx, "handling request", "request_id", requestID)
package request_id
