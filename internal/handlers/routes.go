package handlers

import (
	"net/http"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /api/archive/information", http.HandlerFunc(app.HandleArchiveInformation))
	mux.Handle("POST /api/archive/files", http.HandlerFunc(app.HandleCreateArchive))
	mux.Handle("POST /api/mail/file", http.HandlerFunc(app.HandleSendFile))

	standard := New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}
