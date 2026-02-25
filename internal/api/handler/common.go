package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// logWithTrace enriches a logrus entry with trace context from ctx.
func logWithTrace(logger *logrus.Entry, ctx context.Context) *logrus.Entry {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return logger.WithFields(logrus.Fields{
			"trace_id": spanCtx.TraceID().String(),
			"span_id":  spanCtx.SpanID().String(),
		})
	}
	return logger
}

// isUniqueConstraintError checks if the error is a unique constraint violation
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "unique") ||
		strings.Contains(errStr, "duplicate") ||
		strings.Contains(errStr, "violates unique constraint")
}

// respondError writes an error response with the given HTTP status code.
func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// respondBadRequest writes a 400 Bad Request error response.
func respondBadRequest(c *gin.Context, code, message string) {
	respondError(c, http.StatusBadRequest, code, message)
}

// respondNotFound writes a 404 Not Found error response.
func respondNotFound(c *gin.Context, code, message string) {
	respondError(c, http.StatusNotFound, code, message)
}

// respondInternalError logs the error and writes a 500 Internal Server Error response.
func respondInternalError(c *gin.Context, logger *logrus.Entry, err error, logMsg, responseMsg string) {
	logger.WithError(err).Error(logMsg)
	respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", responseMsg)
}

// respondServiceUnavailable writes a 503 Service Unavailable error response.
func respondServiceUnavailable(c *gin.Context, message string) {
	respondError(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}

// parseUintID parses a uint64 ID from the given gin URL parameter.
// Returns the parsed ID and true on success. On failure it writes a 400 error
// response and returns 0 and false.
func parseUintID(c *gin.Context, paramName string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(paramName), 10, 64)
	if err != nil {
		respondBadRequest(c, "INVALID_ID", "Invalid "+paramName)
		return 0, false
	}
	return id, true
}

// bindJSON binds the request body to the provided struct.
// Returns true on success. On failure it writes a 400 error response and returns false.
func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		respondBadRequest(c, "INVALID_REQUEST", err.Error())
		return false
	}
	return true
}

// requireStringParam validates that a URL parameter is non-empty.
// Returns the value and true on success. On failure it writes a 400 error response
// and returns an empty string and false.
func requireStringParam(c *gin.Context, paramName, errorCode, errorMessage string) (string, bool) {
	val := c.Param(paramName)
	if val == "" {
		respondBadRequest(c, errorCode, errorMessage)
		return "", false
	}
	return val, true
}

// handleDBLookupError handles the common pattern of checking a GORM lookup result
// for ErrRecordNotFound vs other errors. Returns true if an error was handled
// (caller should return). Returns false if there was no error.
func handleDBLookupError(c *gin.Context, logger *logrus.Entry, err error, notFoundCode, notFoundMsg, logMsg, responseMsg string) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondNotFound(c, notFoundCode, notFoundMsg)
		return true
	}
	respondInternalError(c, logger, err, logMsg, responseMsg)
	return true
}

// parsePagination extracts page and page_size query parameters with defaults
// and bounds checking. Returns page, pageSize, and the computed offset.
func parsePagination(c *gin.Context, defaultPageSize, maxPageSize int) (page, pageSize, offset int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(defaultPageSize)))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > maxPageSize {
		pageSize = defaultPageSize
	}

	offset = (page - 1) * pageSize
	return page, pageSize, offset
}

// requireDependency checks that a dependency is not nil and returns true if available.
// On failure it writes a 503 response and returns false.
func requireDependency(c *gin.Context, dep any, name string) bool {
	if dep == nil {
		respondServiceUnavailable(c, name+" not available")
		return false
	}
	return true
}

// parseLimit extracts a limit query parameter with default and bounds checking.
func parseLimit(c *gin.Context, defaultLimit, maxLimit int) int {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if limit < 1 || limit > maxLimit {
		limit = defaultLimit
	}
	return limit
}

// totalPages computes the number of pages given a total count and page size.
func totalPages(total int64, pageSize int) int {
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	return pages
}
