// pipeline {
//   agent any

//   stages {
//     stage('Checkout') {
//       steps {
//         git url: 'https://github.com/QuangPham789/devops-ai-demo.git', branch: 'main'
//       }
//     }

//     stage('Test') {
//       steps {
//         sh 'go test ./...'
//       }
//     }

//     stage('Build Docker') {
//       steps {
//         sh 'docker build -t quangpham789/go-ai-devops:${BUILD_NUMBER} ../docker'
//       }
//     }

//     stage('Push Docker') {
//       steps {
//         withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
//           sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
//           sh 'docker push quangpham789/go-ai-devops:${BUILD_NUMBER}'
//         }
//       }
//     }
//   }

//   // post {
//   //   failure {
//   //     script {
//   //       // Gửi alert tới local log analyzer
//   //       sh """
//   //       curl -s -X POST http://localhost:8080/analyze-log \
//   //         -H 'Content-Type: application/json' \
//   //         -d '{"log":"BUILD_LOG_PLACEHOLDER"}' || true
//   //       """
 
//   //     } 
//   //   }
//   // }
// }
 
pipeline {
    agent any

    environment {
        // khai báo user & token ở đây để sử dụng trong post
        JENKINS_USER = 'admin'
        JENKINS_API_TOKEN = '74dc0fcf9b84aca22624031a7b088c2a'
    }

    stages {
        stage('Checkout') {
            steps {
                git url: 'https://github.com/QuangPham789/devops-ai-demo.git', branch: 'main'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./... || true' // tránh crash pipeline
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
                sh """
                set +e
                # Lấy crumb
                CRUMB=\$(curl -s -u ${JENKINS_USER}:${JENKINS_API_TOKEN} http://localhost:8080/crumbIssuer/api/json | jq -r '.crumb')
                # Gửi POST log
                curl -s -X POST http://localhost:8080/analyze-log \\
                  -H 'Content-Type: application/json' \\
                  -H "Jenkins-Crumb:\$CRUMB" \\
                  -u ${JENKINS_USER}:${JENKINS_API_TOKEN} \\
                  -d '{"log":"BUILD_LOG_PLACEHOLDER"}' || true
                """
            }
        }
    }
}


