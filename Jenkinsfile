pipeline {
    agent any

    environment {
        DOCKERHUB_CREDENTIALS = 'dockerhub'
    }

    stages {
        stage('Checkout') {
            steps {
                // Nếu repo đã clone, dùng pull để tránh lỗi
                sh '''
                if [ -d ".git" ]; then
                    git reset --hard
                    git clean -fd
                    git pull origin main 2>&1 | tee -a build.log
                else
                    git clone -b main https://github.com/QuangPham789/devops-ai-demo.git . 2>&1 | tee -a build.log
                fi
                '''
            }
        }

        stage('Test') {
            steps {
                // Capture log nhưng vẫn fail nếu go test fail
                sh '''
                go test ./... 2>&1 | tee -a build.log
                test ${PIPESTATUS[0]} -eq 0
                '''
            }
        }

        stage('Build Docker') {
            steps {
                sh '''
                docker build -t quangpham789/go-ai-devops:${BUILD_NUMBER} ../docker 2>&1 | tee -a build.log
                test ${PIPESTATUS[0]} -eq 0
                '''
            }
        }

        stage('Push Docker') {
            steps {
                withCredentials([usernamePassword(credentialsId: env.DOCKERHUB_CREDENTIALS, usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh '''
                    echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin 2>&1 | tee -a build.log
                    docker push quangpham789/go-ai-devops:${BUILD_NUMBER} 2>&1 | tee -a build.log
                    test ${PIPESTATUS[1]} -eq 0
                    '''
                }
            }
        }
    }

    post {
        failure {
            script {
                sh '''
                TARGET_URL="https://tennis-scale-tyler-freeze.trycloudflare.com/analyze-log"

                if [ -f build.log ]; then
                    echo "Sending build.log to $TARGET_URL ..."
                    LOG_JSON=$(jq -Rn --arg log "$(cat build.log)" '{log:$log}')
                    curl -v -X POST $TARGET_URL \
                         -H "Content-Type: application/json" \
                         -d "$LOG_JSON" || true
                else
                    echo "build.log not found!"
                fi
                '''
            }
        }
    }
}
