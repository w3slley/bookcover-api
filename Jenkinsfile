pipeline {
    agent any
    stages {
        stage('Checkout code') {
            steps {
                git(url: 'https://github.com/w3slley/bookcover-api.git', branch: 'testing-jenkins')
            }
        }
        stage('Test') {
            steps {
                echo "Testing..."
            }
        }
        stage('Build') {
            steps {
                echo "Building..."
            }
        }
        
    }
}