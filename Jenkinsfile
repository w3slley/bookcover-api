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
        sh 'make test'
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
            https://argocd.infra.longitood.com/api/v1/applications/bookcover-api/resource/actions/v2 \
            -H "Authorization: Bearer $ARGOCD_TOKEN" \
            -H "Content-Type: application/json" \
            -d '{
                "action": "restart",
                "resourceName": "bookcover-api-deployment",
                "kind": "Deployment",
                "group": "apps",
                "version": "v2",
                "namespace": "bookcover-api"
            }'
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
