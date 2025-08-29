pipeline {
  agent any

  environment {
    HARBOR_REGISTRY = 'harbor.infra.longitood.com'
    REPOSITORY_URL = 'git@github.com:w3slley/bookcover-api.git'
    PROD_PROJECT_UUID = 'pw84ocswgoog88sss808gk4c'
    GO_BASE_IMAGE = 'golang:1.23'
    HELM_KUBECTL_BASE_IMAGE = 'baseimages/helm-kubectl'
    HARBOR_CREDENTIALS_ID = 'harbor-jenkins-credentials'
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
            docker.withRegistry("https://${HARBOR_REGISTRY}", "$HARBOR_CREDENTIALS_ID") {
              def customImage = docker.build("bookcover-api/bookcover-api:${env.BRANCH_NAME.replace("/","-")}")
              customImage.push()
            }
        }
      }
    }

    stage('Deploy to ArgoCD') {
      when {
        anyOf {
          branch "main"
        }
      }
      steps {
        withCredentials([string(credentialsId: 'argocd-token', variable: 'ARGOCD_TOKEN')]) {
          sh '''
            curl -X POST \
            https://argocd.infra.longitood.com/api/v1/applications/bookcover-api/sync \
            -H "Authorization: Bearer $ARGOCD_TOKEN" \
            -H "Content-Type: application/json" \
            -d '{"strategy":{"hook":{"force":true}}}'
            '''
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
