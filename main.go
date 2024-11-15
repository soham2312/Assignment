package main

import (
    "log"
    "net/http"
    "retail-pulse/internal/api"
    "retail-pulse/internal/store"
)

func main() {
    // Initialize store manager with CSV file
    storeMgr, err := store.NewStoreManager("StoreMasterAssignment.csv")
    if err != nil {
        log.Fatalf("Failed to initialize store manager: %v", err)
    }

    router := api.SetupRoutes(storeMgr)
    log.Printf("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}