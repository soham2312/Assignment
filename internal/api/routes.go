package api

import (
    "github.com/gorilla/mux"
    "retail-pulse/internal/store"
)

func SetupRoutes(storeMgr *store.StoreManager) *mux.Router {
    router := mux.NewRouter()
    
    handler := NewHandler(storeMgr)
    router.HandleFunc("/api/submit", handler.handleJobSubmission).Methods("POST")
    router.HandleFunc("/api/status", handler.handleJobStatus).Methods("GET")
    
    return router
}
