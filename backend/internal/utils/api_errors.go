package utils

import "errors"

// Business error sentinels (ADM-01). Handlers map them to HTTP status
// codes (mostly 409 for conflicts / 400 for invalid transitions) and a
// user-facing message; internal details never reach the client.
var (
	ErrTeamCodeConflict        = errors.New("team code already exists")
	ErrTeamInUse               = errors.New("team is referenced by matches, standings or favorites")
	ErrGroupNameConflict       = errors.New("group name already exists")
	ErrGroupInUse              = errors.New("group is referenced by teams or matches")
	ErrCityInUse               = errors.New("city is referenced by stadiums or matches")
	ErrStadiumInUse            = errors.New("stadium is referenced by matches")
	ErrMatchInUse              = errors.New("match is referenced by reminders or favorites")
	ErrInvalidStatusTransition = errors.New("invalid match status transition")
	ErrDuplicateCode           = errors.New("code already exists")
	ErrInvalidTime             = errors.New("invalid time format")
)
