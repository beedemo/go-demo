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
        sh "docker-compose -f docker-compose-test.yml run --rm unit"
        sh "docker build -t go-demo ."
        sh "pwd"
      }
    }
    stage("Staging") {
      steps {
        sh "docker-compose -f docker-compose-test-local.yml up -d staging-dep"
        sh 'HOST_IP=localhost docker-compose -f docker-compose-test-local.yml run --rm staging'
      }
    }
    stage("Scan") {
      steps {
        sh "docker run -d --name anchore_cli -v /var/run/docker.sock:/var/run/docker.sock -v /jenkins:/jenkins anchore/cli:latest"
        sh "docker exec anchore_cli anchore feeds sync"
        sh "docker exec anchore_cli anchore analyze --image go-demo"
        sh "docker exec anchore_cli anchore gate --image go-demo --policy /jenkins/go-demo/anchore_policy.txt"
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
      sh "docker-compose -f docker-compose-test-local.yml down"
      sh "docker rm -f anchore_cli"
    }
  }
}