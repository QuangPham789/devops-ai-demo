package main

import (
    "bytes"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "os"
    "time"
)

type Client struct {
    Endpoint string
    APIKey   string
    Client   *http.Client
}

func NewClient(endpoint, apiKey string) *Client {
    if endpoint == "" {
        endpoint = os.Getenv("GEMINI_API_URL")
    }
    if apiKey == "" {
        apiKey = os.Getenv("GEMINI_API_KEY")
    }
    return &Client{
        Endpoint: endpoint,
        APIKey:   apiKey,
        Client: &http.Client{Timeout: 20 * time.Second},
    }
}

// Generate sends a prompt to the configured Gemini endpoint.
// This function assumes the endpoint accepts {"prompt":"..."} and returns plain text
// or JSON where the main text is in "candidates[0].content" or in the response body.
func (c *Client) Generate(prompt string) (string, error) {
    if c.Endpoint == "" || c.APIKey == "" {
        return "", errors.New("Gemini endpoint or API key not set")
    }

    reqBody := map[string]interface{}{"prompt": prompt}
    jb, _ := json.Marshal(reqBody)
    req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(jb))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.APIKey)

    resp, err := c.Client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        b, _ := io.ReadAll(resp.Body)
        return "", errors.New(string(b))
    }
    b, _ := io.ReadAll(resp.Body)
    // Try to parse common JSON shapes
    var out map[string]interface{}
    if err := json.Unmarshal(b, &out); err == nil {
        // Check common fields
        if candidates, ok := out["candidates"].([]interface{}); ok && len(candidates) > 0 {
            if c0, ok := candidates[0].(map[string]interface{}); ok {
                if content, ok := c0["content"].(string); ok {
                    return content, nil
                }
            }
        }
        if text, ok := out["text"].(string); ok {
            return text, nil
        }
        if choices, ok := out["choices"].([]interface{}); ok && len(choices) > 0 {
            if ch0, ok := choices[0].(map[string]interface{}); ok {
                if msg, ok := ch0["message"].(map[string]interface{}); ok {
                    if content, ok := msg["content"].(string); ok {
                        return content, nil
                    }
                }
            }
        }
    }
    // fallback to raw body
    return string(b), nil
}
