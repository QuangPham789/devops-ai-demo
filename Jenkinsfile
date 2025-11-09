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

  // post {
  //   failure {
  //     script {
  //       // Gửi alert tới local log analyzer
  //       sh """
  //       curl -s -X POST http://localhost:8080/analyze-log \
  //         -H 'Content-Type: application/json' \
  //         -d '{"log":"BUILD_LOG_PLACEHOLDER"}' || true
  //       """
 
  //     } 
  //   }
  // }
}
post {
    failure {
        script {
            // Lấy crumb từ Jenkins
            def crumbJson = sh(
                script: "curl -s -u ${JENKINS_USER}:${JENKINS_API_TOKEN} http://localhost:8080/crumbIssuer/api/json",
                returnStdout: true
            ).trim()

            def crumb = sh(
                script: "echo '${crumbJson}' | jq -r '.crumb'",
                returnStdout: true
            ).trim()

            // Gửi alert tới local log analyzer
            sh """
            curl -s -X POST http://localhost:8080/analyze-log \
              -H 'Content-Type: application/json' \
              -H 'Jenkins-Crumb:${crumb}' \
              -u ${JENKINS_USER}:${JENKINS_API_TOKEN} \
              -d '{"log":"BUILD_LOG_PLACEHOLDER"}' || true
            """
        }
    }
}

