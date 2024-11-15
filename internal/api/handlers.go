package api

import (
    "encoding/json"
    "net/http"
    "retail-pulse/internal/models"
    "retail-pulse/internal/processor"
    "retail-pulse/internal/store"
    "sync"
    "time"

    "github.com/google/uuid"
)

type Handler struct {
    jobs     map[string]*models.Job
    jobsMutex sync.RWMutex
    storeMgr  *store.StoreManager
}

func NewHandler(storeMgr *store.StoreManager) *Handler {
    return &Handler{
        jobs:     make(map[string]*models.Job),
        storeMgr: storeMgr,
    }
}

func (h *Handler) handleJobSubmission(w http.ResponseWriter, r *http.Request) {
    var submission models.JobSubmission
    if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate submission
    if submission.Count != len(submission.Visits) {
        http.Error(w, `{ "error": "Count is Wrong"}`, http.StatusBadRequest)
        return
    }

    // Validate store IDs and URLs
    for _, visit := range submission.Visits {
        if err := h.storeMgr.ValidateStore(visit.StoreID); err != nil {
            http.Error(w, "Invalid store ID: "+visit.StoreID, http.StatusBadRequest)
            return
        }
        for _, url := range visit.ImageURLs {
            if url == "" {
                http.Error(w, "Empty image URL for store: "+visit.StoreID, http.StatusBadRequest)
                return
            }
        }
    }

    // Create job
    jobID := uuid.New().String()
    job := &models.Job{
        ID:        jobID,
        Status:    "ongoing",
        Visits:    submission.Visits,
        CreatedAt: time.Now(),
    }

    h.jobsMutex.Lock()
    h.jobs[jobID] = job
    h.jobsMutex.Unlock()

    // Start processing in background
    go h.processJob(job)

    // Return response
    response := models.JobResponse{JobID: jobID}
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleJobStatus(w http.ResponseWriter, r *http.Request) {
    jobID := r.URL.Query().Get("jobid")
    if jobID == "" {
        http.Error(w, "Enter the JOB ID", http.StatusBadRequest)
        return
    }

    h.jobsMutex.RLock()
    job, exists := h.jobs[jobID]
    h.jobsMutex.RUnlock()

    if !exists {
        http.Error(w, "{}", http.StatusNotFound)
        return
    }

    status := models.JobStatus{
        Status: "completed",
        JobID:  job.ID,
    }

    if job.Status == "failed" {
        for _, result := range job.Results {
            if result.Error != nil {
                status.Errors = append(status.Errors, models.JobError{
                    StoreID: result.StoreID,
                    Error:   result.Error.Error(),
                })
            }
        }
    }

    json.NewEncoder(w).Encode(status)
}

func (h *Handler) processJob(job *models.Job) {
    var wg sync.WaitGroup
    results := make(chan models.ImageResult)
    
    // Process each visit's images
    for _, visit := range job.Visits {
        for _, imageURL := range visit.ImageURLs {
            wg.Add(1)
            go func(storeID, url string) {
                defer wg.Done()
                perimeter, err := processor.ProcessImage(url)
                results <- models.ImageResult{
                    StoreID:   storeID,
                    ImageURL:  url,
                    Perimeter: perimeter,
                    Error:     err,
                }
            }(visit.StoreID, imageURL)
        }
    }

    // Close results channel when all goroutines complete
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    var failed bool
    for result := range results {
        job.Results = append(job.Results, result)
        if result.Error != nil {
            failed = true
        }
    }

    h.jobsMutex.Lock()
    if failed {
        job.Status = "failed"
    } else {
        job.Status = "completed"
    }
    h.jobsMutex.Unlock()
}