package handlers

import (
	"project-cahillaw/servers/gateway/models/users"
	"project-cahillaw/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store


type HandlerContext struct {
	Key           string         `json:"key"`
	SessionsStore sessions.Store `json:"sessionsstore"`
	UserStore     users.Store    `json:"userstore"`
}
