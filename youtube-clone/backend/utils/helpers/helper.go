package helpers

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GinContextToContextMiddleware converts Gin context to standard context
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "YOUTUBE_API", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// PaginationParams represents parameters for pagination
type PaginationParams struct {
	Page          int    // Page number
	Size          int    // Page size
	SortBy        string // Field to sort by
	SortDirection string // Sorting direction (asc/desc)
}

// ParsePaginationParams parses pagination parameters from Gin context
func ParsePaginationParams(c *gin.Context) PaginationParams {
	defaultSize := 10
	defaultPage := 1

	size := getIntQueryParam(c, "size", defaultSize)
	page := getIntQueryParam(c, "page", defaultPage)
	sortBy := c.Query("sort_by")
	sortDirection := c.Query("sort_direction")

	return PaginationParams{
		Page:          page,
		Size:          size,
		SortBy:        sortBy,
		SortDirection: sortDirection,
	}
}

// getIntQueryParam extracts integer query parameter from Gin context
func getIntQueryParam(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// fieldError represents a validation error on a specific field
type fieldError struct {
	err validator.FieldError
}

// String returns the string format of the field error
func (fe fieldError) String() string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + fe.err.Field() + "'")
	sb.WriteString(", condition: must be " + fe.err.ActualTag())

	// Print condition parameters, e.g. one of=red blue -> { red blue }
	if fe.err.Param() != "" {
		sb.WriteString(" { " + fe.err.Param() + " }")
	}

	if fe.err.Value() != nil && fe.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", fe.err.Value()))
	}

	return sb.String()
}

// ValidateRequest validates the request struct to ensure it matches requirements
func ValidateRequest(request interface{}) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	err := validate.Struct(request)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			return fmt.Errorf(fieldError{fieldErr}.String())
		}
	}

	return nil
}
