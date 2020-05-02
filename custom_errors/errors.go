package custom_errors

import "errors"

var (
	ErrUserNotOwner                   = errors.New("current user is not owner of this entity, and hence cannot modify it")
	ErrNoEntityMatchingId             = errors.New("no entity found that matches given id")
	ErrOwnerApplyToOwnJob             = errors.New("owner cannot apply to their job")
	ErrApplicationWithdrawnOrRejected = errors.New("owner cannot modify applications with withdrawn status")
	ErrInvalidNewApplicationState     = errors.New("owner cannot move application status to withdrawn or pending")
	ErrJobAlreadyCompleted            = errors.New("job is already completed")
	ErrEntityDeleted                  = errors.New("entity was deleted")
	ErrUserNotAuthenticated           = errors.New("unauthorized request")
	ErrInvalidId                      = errors.New("invalid id supplied")
	ErrInvalidQuery                   = errors.New("invalid query")
)
