package handlers

import (
	"net/http"
)

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /api/archive/information", http.HandlerFunc(h.HandleArchiveInformation))
	mux.Handle("POST /api/archive/files", http.HandlerFunc(h.HandleCreateArchive))
	mux.Handle("POST /api/mail/file", http.HandlerFunc(h.HandleSendFile))

	standard := New(h.recoverPanic, h.logRequest, secureHeaders)
	return standard.Then(mux)
}
