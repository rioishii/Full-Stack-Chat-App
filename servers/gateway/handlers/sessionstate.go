package handlers

import (
	"time"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/models/users"
)

//SessionState represents the user's time at which the session began
//and the authenticated user who started the session
type SessionState struct {
	SessionBegin time.Time
	User         *users.User
}
