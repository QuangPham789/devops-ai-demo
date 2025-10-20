# go-ai-devops

Demo repository: Integrating Gemini (AI) into CI/CD & Monitoring using Go.

## Mục tiêu
- Dùng Gemini API để generate Dockerfile / Jenkinsfile từ prompt
- Khi Jenkins job fail, dịch log và trả summary + suggested fix lên Slack
- Khi Prometheus alert xảy ra, gửi alert qua service → Gemini tạo incident summary → post Slack

## Yêu cầu
- Docker & Docker Compose
- Go 1.20+
- Gemini API key (set trong .env)
- Slack incoming webhook (set trong .env)

## Cấu trúc
- app/: Go web service
- ai/: Gemini client wrapper
- jenkins/Jenkinsfile: sample pipeline (gọi AI trên failure)
- docker/Dockerfile: Dockerfile cho Go app
- monitoring/: prometheus config + alert simulator
- docker-compose.yml: compose to run app + prometheus + jenkins
- slides/: markdown slides + speaker notes
- prompts/: example prompts

## Cách chạy nhanh (local)
1. Sao chép `.env.example` thành `.env` và điền:
   - GEMINI_API_URL (ví dụ: https://YOUR_GEMINI_ENDPOINT)
   - GEMINI_API_KEY
   - SLACK_WEBHOOK_URL

2. Build image app:
   ```
   docker build -t go-ai-devops/app:local -f docker/Dockerfile .
   ```

3. Khởi docker-compose:
   ```
   docker-compose up --build
   ```
   - Go app: http://localhost:8080
   - Jenkins: http://localhost:8081
   - Prometheus: http://localhost:9090

4. Test endpoints:
   - Generate Dockerfile from prompt:
     POST http://localhost:8080/generate/docker with JSON `{"prompt":"Write Dockerfile for Go app ..."}`

   - Analyze log:
     POST http://localhost:8080/analyze-log with JSON `{"log":"..."}`
   - Simulate alert:
     Run `go run monitoring/alert_simulator.go` or use HTTP POST to /alert on app.

## Lưu ý
- Gemini API/URL phụ thuộc vào bạn: mình để GEMINI_API_URL để dễ cấu hình.
- Jenkinsfile mẫu dùng curl để post build log tới service /analyze-log.

## Slide
Slides markdown with speaker notes in `slides/AI_in_DevOps.md`.
