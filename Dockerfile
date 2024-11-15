FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy the entire project including vendor directory
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy only the necessary files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/StoreMasterAssignment.csv .

EXPOSE 8080

CMD ["./main"]