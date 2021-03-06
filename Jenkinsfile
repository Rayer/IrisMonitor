Project_Name = 'Iris-Monitor'
Iris = 'node.rayer.idv.tw'
Iris_OCR1 = 'node1.rayer.idv.tw'
Iris_OCR2 = 'node2.rayer.idv.tw'
pipeline {
    agent any
    parameters {
        string defaultValue: 'master', description: 'Branch name to deploy on server', name: 'branch', trim: false
        string defaultValue: 'monitor-server.app', description: 'Monitor server name', name: 'server_app', trim: false
        string defaultValue: '12800', description: 'Monitor port', name: 'server_port', trim: false
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
                sh label: 'build server', script: "go build -o bin/${params.server_app}"
                archiveArtifacts artifacts: 'bin/*', fingerprint: true, followSymlinks: false, onlyIfSuccessful: true
            }
        }

        stage('Push executable to servers') {
            steps {
                sh label: "Stopping monitor on ${Iris}", script: "ssh jenkins@${Iris} -C pkill ${params.server_app} -f || true"
                sh label: "SCP to ${Iris}", script: "scp bin/${params.server_app} jenkins@${Iris}:~/"
                sh label: "Starting monitor on ${Iris}", script: "ssh jenkins@${Iris} -f sh -c `nohup ./${params.server_app} & `"
                sh label: "Stopping monitor on ${Iris_OCR1}", script: "ssh jenkins@${Iris_OCR1} -C pkill ${params.server_app} -f || true"
                sh label: "SCP to ${Iris_OCR1}", script: "scp bin/${params.server_app} jenkins@${Iris_OCR1}:~/"
                sh label: "Starting monitor on ${Iris_OCR1}", script: "ssh jenkins@${Iris_OCR1} -f sh -c `nohup ./${params.server_app} & `"
                sh label: "Stopping monitor on ${Iris_OCR2}", script: "ssh jenkins@${Iris_OCR2} -C pkill ${params.server_app} -f || true"
                sh label: "SCP to ${Iris_OCR2}", script: "scp bin/${params.server_app} jenkins@${Iris_OCR2}:~/"
                sh label: "Starting monitor on ${Iris_OCR2}", script: "ssh jenkins@${Iris_OCR2} -f sh -c `nohup ./${params.server_app} & `"
            }
        }
    }

   post {
        aborted {
            slackSend message: "Project ${Project_Name} aborted."
        }
        success {
            slackSend message: "Project ${Project_Name} is built successfully."
        }
        failure {
            slackSend message: "Project ${Project_Name} is failed to be built."
        }
    }
}