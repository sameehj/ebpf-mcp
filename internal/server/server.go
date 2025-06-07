package server

import (
	"embed"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sameehj/ebpf-mcp/internal/core"
)

//go:embed .well-known/*
var embeddedFiles embed.FS

func Start() error {
	router := mux.NewRouter()

	// Serve /.well-known/mcp.json from embedded filesystem
	router.PathPrefix("/.well-known/mcp/").Handler(
		http.FileServer(http.FS(embeddedFiles)),
	)

	router.HandleFunc("/rpc", core.HandleMCP).Methods("POST")

	return http.ListenAndServe(":8080", router)
}
