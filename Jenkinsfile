script {
    def short_commit = null
}
pipeline {
  options { 
    buildDiscarder(logRotator(numToKeepStr: '5')) 
    skipDefaultCheckout() 
  }
  agent {
    label "dind-compose"
  }
  environment {
    DOCKER_HUB_USER = 'beedemo'
    DOCKER_CREDENTIAL_ID = 'docker-hub-beedemo'
  }
  stages {
    stage("Unit") {
      steps {
        checkout scm
        gitShortCommit(7)
        sh "docker-compose -f docker-compose-test.yml -p ${BUILD_NUMBER}-${SHORT_COMMIT} run --rm unit"
        script {
          env.IMAGE_ID = sh(returnStdout: true, script: "docker build -q .").trim()
        }
      }
    }
    stage("Staging") {
      steps {
        sh "IMAGE_ID=${IMAGE_ID} docker-compose -f docker-compose-test-local.yml -p ${BUILD_NUMBER}-${SHORT_COMMIT} up -d staging-dep"
        sh "HOST_IP=localhost docker-compose -f docker-compose-test-local.yml -p ${BUILD_NUMBER}-${SHORT_COMMIT} run --rm staging"
      }
    }
    stage("Publish") {
      steps {
        sh "docker tag ${IMAGE_ID} ${DOCKER_HUB_USER}/go-demo:${SHORT_COMMIT}"
        withDockerRegistry(registry: [credentialsId: "$DOCKER_CREDENTIAL_ID"]) {
           sh "docker push $DOCKER_HUB_USER/go-demo:${SHORT_COMMIT}"
        }
      }
    }
  }
  post {
    always {
      sh "docker-compose -f docker-compose-test-local.yml -p ${BUILD_NUMBER}-${SHORT_COMMIT} down"
    }
  }
}