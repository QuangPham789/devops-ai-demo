package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	// "go-ai-devops/ai"
)

type requestPrompt struct {
    Prompt string `json:"prompt"`
}
type logReq struct {
    Log string `json:"log"`
}

var gemini *Client
var slackWebhook string

func init() {
    _ = godotenv.Load()
    gemini = NewClient(os.Getenv("GEMINI_API_URL"), os.Getenv("GEMINI_API_KEY"))
    slackWebhook = os.Getenv("SLACK_WEBHOOK_URL")
}

func main() {
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "go-ai-devops service running")
    })
    http.HandleFunc("/generate/docker", generateDockerHandler)
    http.HandleFunc("/generate/jenkins", generateJenkinsHandler)
    http.HandleFunc("/analyze-log", analyzeLogHandler)
    http.HandleFunc("/alert", alertHandler)

    srv := &http.Server{
        Addr:              ":" + port,
        ReadHeaderTimeout: 5 * time.Second,
    }
    log.Printf("Listening on %s\n", srv.Addr)
    log.Fatal(srv.ListenAndServe())
}

func generateDockerHandler(w http.ResponseWriter, r *http.Request) {
    var req requestPrompt
    body, _ := io.ReadAll(r.Body)
    _ = json.Unmarshal(body, &req)
    if req.Prompt == "" {
        http.Error(w, "prompt required", 400)
        return
    }
    resp, err := gemini.Generate(req.Prompt)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    writeJSON(w, map[string]string{"dockerfile": resp})
}

func generateJenkinsHandler(w http.ResponseWriter, r *http.Request) {
    var req requestPrompt
    body, _ := io.ReadAll(r.Body)
    _ = json.Unmarshal(body, &req)
    if req.Prompt == "" {
        http.Error(w, "prompt required", 400)
        return
    }
    resp, err := gemini.Generate(req.Prompt)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    writeJSON(w, map[string]string{"jenkinsfile": resp})
}

func analyzeLogHandler(w http.ResponseWriter, r *http.Request) {
    var req logReq
    body, _ := io.ReadAll(r.Body)
    _ = json.Unmarshal(body, &req)
    if req.Log == "" {
        http.Error(w, "log required", 400)
        return
    }
    prompt := "Summarize this Jenkins build log and suggest a concise fix. Log:\n\n" + req.Log
    resp, err := gemini.Generate(prompt)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    // post to slack (non-blocking)
    go postToSlack(resp)

    writeJSON(w, map[string]string{"summary": resp})
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
    // Accept generic prometheus alertmanager payload
    body, _ := io.ReadAll(r.Body)
    prompt := "You received this Prometheus alert payload. Write an incident summary and suggest next steps:\n\n" + string(body)
    resp, err := gemini.Generate(prompt)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    go postToSlack(resp)
    writeJSON(w, map[string]string{"incident_summary": resp})
}

func postToSlack(msg string) {
    if slackWebhook == "" {
        log.Println("SLACK_WEBHOOK_URL not set; skipping slack post")
        return
    }
    payload := map[string]string{"text": msg}
    b, _ := json.Marshal(payload)
    http.Post(slackWebhook, "application/json", bytes.NewBuffer(b))
}

func writeJSON(w http.ResponseWriter, v interface{}) {
    b, _ := json.MarshalIndent(v, "", "  ")
    w.Header().Set("Content-Type", "application/json")
    w.Write(b)
}
