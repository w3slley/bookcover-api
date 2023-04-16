pipeline {
    agent any
    stages {
        stage('Checkout code') {
            steps {
                git(url: 'https://github.com/w3slley/bookcover-api.git', branch: 'testing-jenkins')
            }
        }
        stage('Testing') {
            steps {
                sh 'npm install'
                sh 'npm test'
            }
        }
    }
}