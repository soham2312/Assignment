package models

import "time"

type JobSubmission struct {
    Count  int      `json:"count"`
    Visits []Visit  `json:"visits"`
}

type Visit struct {
    StoreID   string   `json:"store_id"`
    ImageURLs []string `json:"image_url"`
    VisitTime string   `json:"visit_time"`
}

type JobResponse struct {
    JobID string `json:"job_id"`
}

type JobStatus struct {
    Status string     `json:"status"`
    JobID  string     `json:"job_id"`
    Errors []JobError `json:"error,omitempty"`
}

type JobError struct {
    StoreID string `json:"store_id"`
    Error   string `json:"error"`
}

type ImageResult struct {
    StoreID   string
    ImageURL  string
    Perimeter float64
    Error     error
}

type Job struct {
    ID        string
    Status    string
    Visits    []Visit
    Results   []ImageResult
    CreatedAt time.Time
}
