package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func redirect(ctx *Context, url string) {
	ctx.Writer.Header().Set("Location", url)
	ctx.Writer.WriteHeader(302)
	ctx.Writer.Write([]byte(url))
}

func checkAuthority(ctx *Context) (*UserPtr, bool) {
	if key, res := GetSessionKey(ctx); res {
		if user, ext := auth.CheckSession(key); ext {
			return user, true
		}
	}
	ClearSession(ctx)
	return nil, false
}

func hPostLogin(ctx *Context) {
	ClearSession(ctx)

	id := ctx.Request.PostFormValue("id")
	pw := ctx.Request.PostFormValue("pw")

	usr, res := auth.Verify(id, pw)
	if !res {
		ctx.Writer.Write([]byte("failed"))
		return
	}

	SetSessionKey(ctx, usr.sessionKey)
	ctx.Writer.Write([]byte("success"))
}

func hAnyLogout(ctx *Context) {
	if key, res := GetSessionKey(ctx); res {
		auth.Logout(key)
	}
	ClearSession(ctx)
	redirect(ctx, "/login")
}

func hGetLogin(ctx *Context) {
	if _, c := checkAuthority(ctx); c {
		redirect(ctx, "/account")
		return
	}
	t, _ := template.ParseFiles("web/templates/login.html")
	t.Execute(ctx.Writer, nil)
}

func hGetAccount(ctx *Context) {
	if user, c := checkAuthority(ctx); c {
		ctx.Writer.Write([]byte(user.user.name))
		return
	}
	redirect(ctx, "/login")
}

func hGetStaticFile(ctx *Context) {
	localPath := filepath.Join("web/static/", ctx.Params["path"])
	content, err := ioutil.ReadFile(localPath)
	if err != nil {
		ctx.Writer.WriteHeader(404)
		ctx.Writer.Write([]byte(http.StatusText(404)))
		return
	}

	contentType := getContentType(&content, localPath)
	ctx.Writer.Header().Add("Content-Type", contentType)
	ctx.Writer.Write(content)
}

func getContentType(buffer *[]byte, path string) string {
	ext := filepath.Ext(path)
	var contentType string
	switch ext {
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpeg"
	default:
		contentType = http.DetectContentType(*buffer)
	}

	return contentType
}
