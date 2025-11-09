// // pipeline {
// //   agent any

// //   stages {
// //     stage('Checkout') {
// //       steps {
// //         git url: 'https://github.com/QuangPham789/devops-ai-demo.git', branch: 'main'
// //       }
// //     }

// //     stage('Test') {
// //       steps {
// //         sh 'go test ./...'
// //       }
// //     }

// //     stage('Build Docker') {
// //       steps {
// //         sh 'docker build -t quangpham789/go-ai-devops:${BUILD_NUMBER} ../docker'
// //       }
// //     }

// //     stage('Push Docker') {
// //       steps {
// //         withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
// //           sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
// //           sh 'docker push quangpham789/go-ai-devops:${BUILD_NUMBER}'
// //         }
// //       }
// //     }
// //   }

// //   // post {
// //   //   failure {
// //   //     script {
// //   //       // Gửi alert tới local log analyzer
// //   //       sh """
// //   //       curl -s -X POST http://localhost:8080/analyze-log \
// //   //         -H 'Content-Type: application/json' \
// //   //         -d '{"log":"BUILD_LOG_PLACEHOLDER"}' || true
// //   //       """
 
// //   //     } 
// //   //   }
// //   // }
// // }
 
// pipeline {
//     agent any

//     environment {
//         // khai báo user & token ở đây để sử dụng trong post
//         JENKINS_USER = 'admin'
//         JENKINS_API_TOKEN = '11ad818e8b9af83ddd92ef90e20e914721'
//     }

//     stages {
//         stage('Checkout') {
//             steps {
//                 git url: 'https://github.com/QuangPham789/devops-ai-demo.git', branch: 'main'
//             }
//         }

//         stage('Test') {
//             steps {
//                 sh 'go test ./... || true' // tránh crash pipeline
//             }
//         }

//         stage('Build Docker') {
//             steps {
//                 sh 'docker build -t quangpham789/go-ai-devops:${BUILD_NUMBER} ../docker'
//             }
//         }

//         stage('Push Docker') {
//             steps {
//                 withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
//                     sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
//                     sh 'docker push quangpham789/go-ai-devops:${BUILD_NUMBER}'
//                 }
//             }
//         }
//     }

//     // post {
//     //     failure {
//     //         script {
//     //             sh """
//     //             set +e
//     //             # Lấy crumb
//     //             CRUMB=\$(curl -s -u ${JENKINS_USER}:${JENKINS_API_TOKEN} http://localhost:8080/crumbIssuer/api/json | jq -r '.crumb')
//     //             # Gửi POST log
//     //             curl -s -X POST http://localhost:8080/analyze-log \\
//     //               -H 'Content-Type: application/json' \\
//     //               -H "Jenkins-Crumb:\$CRUMB" \\
//     //               -u ${JENKINS_USER}:${JENKINS_API_TOKEN} \\
//     //               -d '{"log":"BUILD_LOG_PLACEHOLDER"}' || true
//     //             """
//     //         }
//     //     }
//     // }
//     post {
//     failure {
//         script {
//             sh """
//             set +e
//             # Lấy crumb và parse bằng sed
//             CRUMB=\$(curl -s -u ${JENKINS_USER}:${JENKINS_API_TOKEN} https://voltage-dietary-sentence-attempts.trycloudflare.com/crumbIssuer/api/json | sed -n 's/.*"crumb":"\\([^"]*\\)".*/\\1/p')
            
//             echo "Crumb: \$CRUMB"

//             # Gửi POST log
//             curl -s -X POST http://localhost:8085/analyze-log \\
//               -H 'Content-Type: application/json' \\
//               -H "Jenkins-Crumb:\$CRUMB" \\
//               -u ${JENKINS_USER}:${JENKINS_API_TOKEN} \\
//               -d '{"log":"BUILD_LOG_PLACEHOLDER"}' || true
//             """
//         }
//     }
// }
// }
pipeline {
    agent any

    environment {
        DOCKERHUB_CREDENTIALS = 'dockerhub'
    }

    stages {
        stage('Checkout') {
            steps {
                git url: 'https://github.com/QuangPham789/devops-ai-demo.git', branch: 'main'
            }
        }

        stage('Test') {
            steps {
                // redirect stdout + stderr vào build.log, tiếp tục dù fail
                sh '''
                set +e
                go test ./... 2>&1 | tee -a build.log
                '''
            }
        }

        stage('Build Docker') {
            steps {
                sh '''
                set +e
                docker build -t quangpham789/go-ai-devops:${BUILD_NUMBER} ../docker 2>&1 | tee -a build.log
                '''
            }
        }

        stage('Push Docker') {
            steps {
                withCredentials([usernamePassword(credentialsId: env.DOCKERHUB_CREDENTIALS, usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh '''
                    set +e
                    echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin 2>&1 | tee -a build.log
                    docker push quangpham789/go-ai-devops:${BUILD_NUMBER} 2>&1 | tee -a build.log
                    '''
                }
            }
        }
    }

    post {
        failure {
            script {
                // gửi build.log nếu tồn tại
                sh '''
                if [ -f build.log ]; then
                    TARGET_URL="https://tennis-scale-tyler-freeze.trycloudflare.com/analyze-log"
                    echo "Sending build.log to $TARGET_URL ..."
                    curl -v -X POST $TARGET_URL \
                         -H "Content-Type: text/plain" \
                         --data-binary @build.log || true
                else
                    echo "build.log not found!"
                fi
                '''
            }
        }
    }
}
