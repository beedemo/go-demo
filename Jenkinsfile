pipeline {
  options { 
    //only keep logs for 5 runs
    buildDiscarder(logRotator(numToKeepStr: '5')) 
    //we only want to ever checkout
    skipDefaultCheckout()
  }
  agent {
    //we need Docker Compose to run many of the sh steps
    //we want to use Docker-in-Docker (DIND) to isolate Docker Compose networks from other running jobs
    label "dind-compose"
  }
  environment {
    //these will be used throughout the Pipeline
    DOCKER_HUB_USER = 'beedemo'
    DOCKER_CREDENTIAL_ID = 'docker-hub-beedemo'
    //will shorten sh step for frist two stages, but require stage level variables to override
    COMPOSE_FILE = 'docker-compose-test.yml'
  }
  stages {
    stage("Prepare Build Environment") {
      steps {
        checkout scm
        //load image in saved in agent
        sh 'docker load -i /jenkins/go-demo-unit-cache.tar'
      }
    }
    stage("Build Cache Image") {
      when {
        //only execute this stage when there are pushes to the build-cache-image branch
        //this makes it easy to update this special, custom Docker image that con
        branch 'build-cache-image'
      }
      steps {
        //this sh step will actually build a customized container using the Dockerfile.build file
        //the docker-compose run command runs a one time command against the specified service
        //the --name option assigns the specified name to the container
        sh "COMMIT_SHA=\$(git rev-parse HEAD | tr -d '\n') docker-compose run --name go-demo-unit unit-cache"
        //now that we have a container running with all of the build dependencies in the container we want to create a new Docker image from it
        //the docker commit command allows us to do exactly that by creating a new image from the go-demo-unit container’s changes
        sh "docker commit go-demo-unit ${DOCKER_HUB_USER}/go-demo:unit-cache"
        //no need to keep the container running, so we will remove it
        sh "docker rm go-demo-unit"
        //sign in to registry - this is the less Declarative way, but results in fewer steps
        withDockerRegistry(registry: [credentialsId: "$DOCKER_CREDENTIAL_ID"]) { 
          //push go-demo specific build image to a Docker registry (DockerHub in this case)
          sh "docker push ${DOCKER_HUB_USER}/go-demo:unit-cache"
        }
      }
    }
    stage("Unit") {
      steps {
        //this global library will sent a SHORT_COMMIT environmental variable the first 7 characters of the commit sha for the current go-demo repo checked out HEAD
        gitShortCommit(7)
        //NOTE: We are using the image that was pushed in the Build Cache Image stage - so if that did not get pushed successfully then this stage will fail
        //the unit service maps the current workspace directory on the dind-compose agent to the '/usr/src/myapp' directory of the unit service container
        //this results in the go-demo binary being created in the workspace and being available for the docker build below
        sh "UNIT_CACHE_IMAGE=${DOCKER_HUB_USER}/go-demo:unit-cache docker-compose run --rm unit"
        junit 'report.xml'
        script {
          //we put this step in a script block - allowing us to fall back to Scripted Pipeline - and in this case assign the output of a sh step to an environmental variable
          //this docker build command uses the "-q" argument which tells it to "Suppress the build output and print image ID on success" - it then strips off the newline character
          //this step diverges from the build step in the DevOps Toolkit 2.1 workshop [sh "docker build -t go-demo ."]
          //this will allow you to avoid naming collisions, although isn't absolutely necessary with DIND vs mounting the Docker Socket where it would be critical in this example
          env.IMAGE_ID = sh(returnStdout: true, script: "docker build --cache-from alpine:3.4 --build-arg COMMIT_SHA=\$(git rev-parse HEAD | tr -d '\n') --build-arg BUILD_CACHE_COMMIT_SHA=\$(docker inspect -f \\\"{{.Config.Labels.commit_sha}}\\\" ${DOCKER_HUB_USER}/go-demo:unit-cache | tr -d '\n') -q . | tr -d '\n'")
        }
      }
    }
    stage("Staging") {
      //do not execute this stage for the build-cache-image branch
      when {
        not { branch 'build-cache-image' }
      }
      steps {
        //we are passing in the ID of the Docker image we built above to use in the compose file
        sh "IMAGE_ID=${IMAGE_ID} docker-compose -f docker-compose-test-local.yml  up -d staging-dep"
        sh "UNIT_CACHE_IMAGE=${DOCKER_HUB_USER}/go-demo:unit-cache HOST_IP=localhost docker-compose -f docker-compose-test-local.yml run --rm staging"
        junit 'report.xml'
      }
    }
    stage("Publish") {
      when {
        branch 'master'
      }
      steps {
        //Note the use of the SHORT_COMMIT environmental variable as the image tag
        sh "docker tag ${IMAGE_ID} ${DOCKER_HUB_USER}/go-demo:${SHORT_COMMIT}"
        //once again setting up credentials for Docker Hub
        withDockerRegistry(registry: [credentialsId: "$DOCKER_CREDENTIAL_ID"]) {
          //push the actual go-app that we just built and tested
          sh "docker push $DOCKER_HUB_USER/go-demo:${SHORT_COMMIT}"
        }
      }
    }
  }
  //the post section allows us to greatly reduce the need for try/catch blocks
  //with Scripted Pipeline we would have had to put a try catch around the entire script
  post {
    always {
      sh "docker-compose -f docker-compose-test-local.yml down"
    }
  }
}
