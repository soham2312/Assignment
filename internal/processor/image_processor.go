package processor

import (
    "errors"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "math/rand"
    "net/http"
    "time"
)

func ProcessImage(url string) (float64, error) {
    if url == "" {
        return 0, errors.New("empty URL")
    }

    // Download image
    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return 0, errors.New("failed to download image")
    }

    // Decode image
    img, _, err := image.Decode(resp.Body)
    if err != nil {
        return 0, err
    }

    bounds := img.Bounds()
    height := bounds.Max.Y - bounds.Min.Y
    width := bounds.Max.X - bounds.Min.X

    // Calculate perimeter
    perimeter := float64(2 * (height + width))

    // Random sleep to simulate processing
    sleepTime := 100 + rand.Intn(301) // 100-400ms
    time.Sleep(time.Duration(sleepTime) * time.Millisecond)

    return perimeter, nil
}