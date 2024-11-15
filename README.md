# Retail Store Image Processing Service

## Description

This service processes images collected from retail stores, calculating image perimeters with simulated GPU processing time. It provides a REST API for job submission and status tracking, handling concurrent processing of multiple jobs with thousands of images.

## Key Features

- REST API endpoints for job submission and status checking
- Concurrent image processing using goroutines
- Store validation against CSV master data
- Image perimeter calculation with simulated processing time
- Job status tracking and error handling
- Docker support for easy deployment

## Assumptions

- Store master data is provided in CSV format with columns: AreaCode, StoreName, StoreID
- Image URLs are publicly accessible
- Supported image formats: JPEG and PNG
- Processing time simulation (0.1-0.4s) is sufficient for demo purposes
- In-memory storage is acceptable for job tracking (not persistent across restarts)
- All images are valid and have readable dimensions
- Network connectivity is available for downloading images

## Installation & Setup

### Prerequisites
- Go 1.21 or later
- Docker
- Git

### Local Setup

#### Clone the repository
```
git clone https://github.com/soham2312/Assignment.git
cd Assignment
```

#### Install dependencies

```
go mod download
```

#### Run the service

```
go run main.go
```

### Docker Setup

```
docker-compose up --build
```

### API Endpoints

#### Submit Job
Create a job to process the images collected from stores.
URL: /api/submit/
Method: POST
Request Payload
```
{
   "count":2,
   "visits":[
      {
         "store_id":"S00339218",
         "image_url":[
            "https://www.gstatic.com/webp/gallery/2.jpg",
            "https://www.gstatic.com/webp/gallery/3.jpg"
         ],
         "visit_time": "time of store visit"
      },
      {
         "store_id":"S01408764",
         "image_url":[
            "https://www.gstatic.com/webp/gallery/3.jpg"
         ],
         "visit_time": "time of store visit"
      }
   ]
}
```

#### Get Job Info
URL : /api/status?jobid=123
URL Parameters: - jobid Job ID received while creating the job
Method: GET
Success Response
Condition: If everything is OK and jobID exists.
Code: 200 OK


## Development Environment
### Hardware/OS

- MacBook Pro M1 (or your actual setup)
- macOS Ventura 13.4 (or your OS)
- 16GB RAM

### Tools

- Visual Studio Code 1.86.0
- Docker Desktop 4.27.0
- Go 1.21.6

### Libraries

- gorilla/mux: HTTP router
- google/uuid: UUID generation
- standard library packages for image processing

## Future Improvements

- Replace in-memory storage with PostgreSQL database for persisting jobs and store master data
- Implement a worker pool with Redis queue for better job processing and scalability
- Add comprehensive monitoring using Prometheus metrics and structured logging
- Set up automated testing pipeline with unit and integration tests
- Add authentication and rate limiting for API security


