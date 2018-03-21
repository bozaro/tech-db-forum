/*
sudo apt install ghp-import
*/
pipeline {
    agent any

    tools {
        go 'go-1.10'
    }

    options {
        skipDefaultCheckout true
    }

    parameters {
        string(name: 'TAG_NAME', defaultValue: '', description: 'Tag name')
    }

    environment {
        GOPATH = "$WORKSPACE"
        PROJ = "github.com/bozaro/tech-db-forum"
        PATH = "$GOPATH/bin:$PATH"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout([
                        $class           : 'GitSCM',
                        branches         : scm.branches,
                        extensions       : scm.extensions + [[$class: 'LocalBranch'], [$class: 'CleanCheckout'], [$class: 'RelativeTargetDirectory', relativeTargetDir: "src/${PROJ}"]],
                        userRemoteConfigs: scm.userRemoteConfigs
                ])
            }
        }
        stage('Prepare') {
            steps {
                dir("src/${PROJ}") {
                    sh """
go install -v ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
go install -v ./vendor/github.com/jteeuwen/go-bindata/go-bindata
go install -v ./vendor/github.com/mailru/easyjson/easyjson
go install -v ./vendor/github.com/aktau/github-release
go generate -x .
"""
                }
                dir("src/${PROJ}/target/dist") {
                    sh "true"
                }
            }
        }
        stage('Build') {
            failFast true
            parallel {
                stage('darwin_amd64') {
                    environment {
                        GOOS = "darwin"
                        GOARCH = "amd64"
                        SUFFIX = ""
                    }
                    steps {
                        dir("src/${PROJ}") {
                            sh """
go build -ldflags " -X ${PROJ}/tests.BuildTag=\${BUILD_TAG} -X ${PROJ}/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                        }
                    }
                }
                stage('linux_386') {
                    environment {
                        GOOS = "linux"
                        GOARCH = "386"
                        SUFFIX = ""
                    }
                    steps {
                        dir("src/${PROJ}") {
                            sh """
go build -ldflags " -X ${PROJ}/tests.BuildTag=\${BUILD_TAG} -X ${PROJ}/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                        }
                    }
                }
                stage('linux_amd64') {
                    environment {
                        GOOS = "linux"
                        GOARCH = "amd64"
                        SUFFIX = ""
                    }
                    steps {
                        dir("src/${PROJ}") {
                            sh """
go build -ldflags " -X ${PROJ}/tests.BuildTag=\${BUILD_TAG} -X ${PROJ}/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                        }
                    }
                }
                stage('windows_386') {
                    environment {
                        GOOS = "windows"
                        GOARCH = "386"
                        SUFFIX = ".exe"
                    }
                    steps {
                        dir("src/${PROJ}") {
                            sh """
go build -ldflags " -X ${PROJ}/tests.BuildTag=\${BUILD_TAG} -X ${PROJ}/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                        }
                    }
                }
                stage('windows_amd64') {
                    environment {
                        GOOS = "windows"
                        GOARCH = "amd64"
                        SUFFIX = ".exe"
                    }
                    steps {
                        dir("src/${PROJ}") {
                            sh """
go build -ldflags " -X ${PROJ}/tests.BuildTag=\${BUILD_TAG} -X ${PROJ}/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                        }
                    }
                }
            }
        }
        stage('Prepare gp-pages') {
            steps {
                dir("src/${PROJ}") {
                    sh """
git branch -fD gh-pages || true
git branch -rd origin/gh-pages || true
ghp-import -n target/dist
"""
                }
            }
        }
        stage('Publish gh-pages') {
            when {
                branch 'master'
            }
            steps {
                withCredentials([usernamePassword(credentialsId: '88e000b8-d989-4f94-b919-1cc1352a5f96', passwordVariable: 'TOKEN', usernameVariable: 'LOGIN')]) {
                    dir("src/${PROJ}") {
                        sh """
git push -qf https://\${TOKEN}@github.com/bozaro/tech-db-forum.git gh-pages
"""
                    }
                }
            }
        }
        stage('Publish release') {
            when {
                expression { params.TAG_NAME != "" }
            }
            environment {
                GITHUB_USER = "bozaro"
                GITHUB_REPO = "tech-db-forum"
            }
            steps {
                withCredentials([[$class: 'StringBinding', credentialsId: '49bf22be-f4d4-4a75-855a-b0e56e357f1c', variable: 'GITHUB_TOKEN']]) {
                    dir("src/${PROJ}") {
                        sh """
github-release info --tag ${params.TAG_NAME} || github-release release --tag ${params.TAG_NAME} --draft
for i in target/dist/*.zip; do
  github-release upload --tag ${params.TAG_NAME} --file \$i --name `basename \$i`
done
"""
                    }
                }
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: "src/${PROJ}/target/dist/*.zip", fingerprint: true
        }
    }
}
