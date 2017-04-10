pipeline {
  options { 
    buildDiscarder(logRotator(numToKeepStr: '5')) 
    skipDefaultCheckout() 
  }
  agent {
    label "docker-compose"
  }
  stages {
    stage("Unit") {
      steps {
        checkout scm
        sh "docker-compose -f docker-compose-test.yml run --rm unit"
        sh "docker build -t go-demo ."
      }
    }
    stage("Staging") {
      steps {
        sh "docker-compose -f docker-compose-test-local.yml up -d staging-dep"
        sh 'HOST_IP=localhost docker-compose -f docker-compose-test-local.yml run --rm staging'
      }
    }
  }
  post {
    always {
      sh "docker-compose -f docker-compose-test-local.yml down"
    }
  }
}