Project_Name = 'Iris-Monitor'
pipeline {
    agent any
    parameters {
        string defaultValue: 'master', description: 'Branch name to deploy on server', name: 'branch', trim: false
        string defaultValue: 'monitor-server.app', description: 'Monitor server name', name: 'server_app', trim: false

    }

    stages {
        stage('Fetch from github') {
            steps {
                slackSend message: "Project ${Project_Name} start to be built"
                git credentialsId: '26c5c0a0-d02d-4d77-af28-761ffb97c5cc', url: 'https://github.com/Rayer/IrisMonitor.git', branch: "${params.branch}"
            }
        }
        stage('Unit test') {
            steps {
                sh label: 'go version', script: 'go version'
                sh label: 'install gocover-cobertura', script: 'go get github.com/t-yuki/gocover-cobertura'
                sh label: 'go unit test', script: 'go test --coverprofile=cover.out'
                sh label: 'convert coverage xml', script: '~/go/bin/gocover-cobertura < cover.out > coverage.xml'
            }
        }
        stage ("Extract test results") {
            steps {
                cobertura coberturaReportFile: 'coverage.xml'
            }
        }

        stage('build and archive executable') {
            steps {
                sh label: 'show version', script: 'go version'
                sh label: 'build server', script: "go build -o ../bin/${params.server_app}"
                archiveArtifacts artifacts: 'bin/*', fingerprint: true, followSymlinks: false, onlyIfSuccessful: true
            }
        }

        stage('Push executable to servers') {
            steps {
                sh label: 'Installing on node.rayer.idv.tw', script: ''
                sh label: 'Installing on node1.rayer.idv.tw', script: ''
                sh label: 'Installing on node2.rayer.idv.tw', script: ''
            }
        }
    }
}