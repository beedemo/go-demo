script {
    def short_commit = null
}
pipeline {
  options { 
    buildDiscarder(logRotator(numToKeepStr: '5')) 
    skipDefaultCheckout() 
  }
  agent {
    label "docker-compose"
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
        sh "docker build -t go-demo ."
      }
    }
    stage("Staging") {
      steps {
        sh "docker-compose -f docker-compose-test-local.yml -p ${BUILD_NUMBER}-${SHORT_COMMIT} up -d staging-dep"
        sh "HOST_IP=localhost docker-compose -f docker-compose-test-local.yml -p ${BUILD_NUMBER}-${SHORT_COMMIT} run --rm staging"
      }
    }
    stage("Publish") {
      steps {
        sh "docker tag go-demo $DOCKER_HUB_USER/go-demo:${SHORT_COMMIT}"
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