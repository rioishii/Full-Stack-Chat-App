package handlers

import (
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/indexes"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/models/users"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-rioishii/servers/gateway/sessions"
)

//HandlerCtx provides access to context for HTTP handler functions
type HandlerCtx struct {
	SigningKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	Trie         *indexes.Trie
	Notifier     *Notifier
}

//NewHandlerContext constructs a new HandlerCtx,
//ensuring that the dependencies are valid values
func NewHandlerContext(signingKey string, sessionStore sessions.Store, userStore users.Store, trie *indexes.Trie, notifier *Notifier) *HandlerCtx {
	if len(signingKey) == 0 {
		panic("nil signing key")
	}
	if sessionStore == nil {
		panic("nil session store")
	}
	if userStore == nil {
		panic("nil user store")
	}
	if trie.Root == nil || trie.Size != 0 {
		panic("nil trie")
	}
	return &HandlerCtx{signingKey, sessionStore, userStore, trie, notifier}
}
