package main

import (
	"github.com/gorilla/sessions"
)

var (
	key   = []byte("ade396da47f808005372837d10ff02231f69b68eb2ebe9139ae2b4a121be7257")
	store = sessions.NewCookieStore(key)
)

//GetSessionKey ...
func GetSessionKey(ctx *Context) (string, bool) {
	session, _ := store.Get(ctx.Request, "session-auth")
	if key, ok := session.Values[SessionKeyUser].(string); ok && key != "" {
		return key, true
	}

	return "", false
}

//SetSessionKey ...
func SetSessionKey(ctx *Context, key string) error {
	session, _ := store.Get(ctx.Request, "session-auth")
	session.Values[SessionKeyUser] = key
	if err := session.Save(ctx.Request, ctx.Writer); err != nil {
		return err
	}
	return nil
}

//ClearSession ...
func ClearSession(ctx *Context) {
	session, _ := store.Get(ctx.Request, "session-auth")
	for k := range session.Values {
		delete(session.Values, k)
	}
	session.Save(ctx.Request, ctx.Writer)
}
