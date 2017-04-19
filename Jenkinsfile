script {
    def short_commit = null
}
pipeline {
  options { 
    buildDiscarder(logRotator(numToKeepStr: '5')) 
    skipDefaultCheckout() 
    timeout(time: 5, unit: 'MINUTES')
  }
  agent {
    label "dind-compose"
  }
  environment {
    DOCKER_HUB_USER = 'beedemo'
    DOCKER_CREDENTIAL_ID = 'docker-hub-beedemo'
  }
  stages {
    stage("Build Cache Image") {
      when {
        branch 'build-cache-image'
      }
      steps {
        checkout scm
        gitShortCommit(7)
        sh "docker-compose -f docker-compose-test.yml -p ${BUILD_NUMBER}-${SHORT_COMMIT} run unit-cache"
        sh "docker ps -a"
        sh "docker commit go-demo-unit ${DOCKER_HUB_USER}/go-demo:unit-cache"
        sh "docker rm go-demo-unit"
        //sign in to registry
        withDockerRegistry(registry: [credentialsId: "$DOCKER_CREDENTIAL_ID"]) { 
            //push repo specific image to Docker registry (DockerHub in this case)
            sh "docker push ${DOCKER_HUB_USER}/go-demo:unit-cache"
        }
      }
    }
    stage("Unit") {
      steps {
        checkout scm
        gitShortCommit(7)
        sh "UNIT_CACHE_IMAGE=${DOCKER_HUB_USER}/go-demo:unit-cache docker-compose -f docker-compose-test.yml -p ${BUILD_NUMBER}-${SHORT_COMMIT} run --rm unit"
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