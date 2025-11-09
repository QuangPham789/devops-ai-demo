pipeline {
  agent any

  stages {
    stage('Checkout') {
      steps {
        git url: 'https://github.com/QuangPham789/devops-ai-demo.git', branch: 'main'
      }
    }

    stage('Test') {
      steps {
        sh 'go test ./...'
      }
    }

    stage('Build Docker') {
      steps {
        sh 'docker build -t quangpham789/go-ai-devops:${BUILD_NUMBER} ../docker'
      }
    }

    stage('Push Docker') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
          sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
          sh 'docker push quangpham789/go-ai-devops:${BUILD_NUMBER}'
        }
      }
    }
  }

  post {
    failure {
      script {
        // Gửi alert tới local log analyzer
        sh """
        curl -s -X POST https://outputs-watts-visual-relay.trycloudflare.com/analyze-log \
          -H 'Content-Type: application/json' \
          -d '{"log":"BUILD_LOG_PLACEHOLDER"}' || true
        """
 
      } 
    }
  }
}
