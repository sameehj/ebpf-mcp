package server

import (
    "net/http"
    "github.com/gorilla/mux"
    "ebpf-mcp/internal/core"
)

func Start() error {
    router := mux.NewRouter()
    router.HandleFunc("/rpc", core.HandleMCP).Methods("POST")

    return http.ListenAndServe(":8080", router)
}
