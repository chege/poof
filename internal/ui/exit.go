package ui

import (
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgconn"
)

// Exit codes following POSIX conventions where applicable.
const (
	ExitOK         = 0
	ExitError      = 1 // General error
	ExitConfigErr  = 2 // Configuration or validation error
	ExitConnErr    = 3 // Database connection error
	ExitPartialErr = 4 // Masking job finished but with row failures
	ExitSafetyErr  = 5 // Safety check failed (unauthorized DB)
)

var (
	// ErrConfig indicates a configuration or validation error.
	ErrConfig = errors.New("configuration error")
	// ErrConnection indicates a database connection failure.
	ErrConnection = errors.New("connection error")
	// ErrSafety indicates a safety check failure.
	ErrSafety = errors.New("safety error")
	// ErrPartial indicates that the job completed but some rows failed to mask.
	ErrPartial = errors.New("masking partial failure")
)

// HandleExit maps an error to a specific exit code and terminates the process.
func HandleExit(err error) {
	if err == nil {
		os.Exit(ExitOK)
	}

	code := ExitError

	// Check for wrapped custom errors
	if errors.Is(err, ErrConfig) {
		code = ExitConfigErr
	} else if errors.Is(err, ErrConnection) {
		code = ExitConnErr
	} else if errors.Is(err, ErrSafety) {
		code = ExitSafetyErr
	} else if errors.Is(err, ErrPartial) {
		code = ExitPartialErr
	}

	// Also check for database-specific connection errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// You could add more specific code mappings here if needed
		code = ExitConnErr
	}

	Error("%s", err.Error())
	os.Exit(code)
}

// WrapError creates a new error wrapping one of our base error types.
func WrapError(base, original error) error {
	if original == nil {
		return base
	}
	return fmt.Errorf("%w: %v", base, original)
}
