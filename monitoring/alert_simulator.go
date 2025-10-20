package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "os"
)

func main() {
    url := os.Getenv("APP_ALERT_ENDPOINT")
    if url == "" {
        url = "http://localhost:8080/alert"
    }
    payload := map[string]interface{}{
        "status": "firing",
        "alerts": []map[string]interface{}{
            {
                "labels": map[string]string{
                    "alertname": "HighCPU",
                    "job":       "go-app",
                    "severity":  "critical",
                },
                "annotations": map[string]string{
                    "summary": "CPU > 90% for 2 minutes",
                },
            },
        },
    }
    b, _ := json.Marshal(payload)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    log.Println("Alert posted, status:", resp.Status)
}
