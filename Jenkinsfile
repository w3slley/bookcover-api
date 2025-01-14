pipeline {
  agent any

  environment {
    DOCKER_REGISTRY = 'https://harbor.infra.longitood.com'
    REPOSITORY_URL = 'git@github.com:w3slley/bookcover-api.git'
    PROD_PROJECT_UUID = 'pw84ocswgoog88sss808gk4c'
    GO_BASE_IMAGE = 'golang:1.23'
  }

  stages {
    stage('Checkout') {
      steps {
        git url: "$REPOSITORY_URL", branch: env.BRANCH_NAME
      }
    }

    stage('Run tests') {
      agent {
        docker {
          image "$GO_BASE_IMAGE"
          args '-e GOCACHE=/tmp/go-build'
        } 
      }
      steps {
        sh 'go test ./tests'
      }
    }
    stage('Build docker image and publish to registry') {
      when {
        anyOf {
          branch "main"
        }
      }
      steps {
        script {
            docker.withRegistry("${DOCKER_REGISTRY}", 'harbor-jenkins-credentials') {
              def customImage = docker.build("bookcover-api/bookcover-api:${env.BRANCH_NAME.replace("/","-")}")
              customImage.push()
            }
        }
      }
    }

    stage('Deploy to coolify') {
      when {
        anyOf {
          branch "main"
        }
      }
      steps {
        script {
          def  uuid = "$PROD_PROJECT_UUID"
          withCredentials([string(credentialsId: 'coolify-token', variable: 'COOLIFY_TOKEN')]) {
            def response = httpRequest(
                url: "https://coolify.longitood.com/api/v1/deploy?uuid=${uuid}",
                httpMode: 'POST',
                contentType: 'APPLICATION_JSON',
                customHeaders: [[name: 'Authorization', value: "Bearer ${COOLIFY_TOKEN}"]],
                requestBody: """
                """
                )
              println("Status: ${response.status}")
              println("Response: ${response.content}")
          }
        }
      }
    }
  }

  post {
    success {
      echo 'Build successful!'
    }
    failure {
      echo 'Build failed!'
    }
  }
}
