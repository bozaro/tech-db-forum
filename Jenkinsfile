pipeline {
    agent {
        dockerfile {
            filename 'Dockerfile.build'
            customWorkspace "/var/jenkins_home/go/src/github.com/bozaro/tech-db-forum"
        }
    }

    parameters {
        string(name: 'TAG_NAME', defaultValue: '', description: 'Tag name')
    }

    environment {
        HOME = "/var/jenkins_home"
        PROJ = "github.com/bozaro/tech-db-forum"
        GOPATH = "/var/jenkins_home/go"
        PATH = "$GOPATH/bin:$PATH"
    }

    stages {
        stage('Prepare') {
            steps {
                sh """
env | sort
export PATH=\$GOPATH/bin:\$PATH
go install -v ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
go install -v ./vendor/github.com/jteeuwen/go-bindata/go-bindata
go install -v ./vendor/github.com/mailru/easyjson/easyjson
go install -v ./vendor/github.com/aktau/github-release
go generate -x .
mkdir -p target/dist
"""
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
                        sh """
export PATH=\$GOPATH/bin:\$PATH
go build -ldflags " -X \$PROJ/tests.BuildTag=\${BUILD_TAG} -X \$PROJ/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                    }
                }
                stage('linux_amd64') {
                    environment {
                        GOOS = "linux"
                        GOARCH = "amd64"
                        SUFFIX = ""
                    }
                    steps {
                        sh """
export PATH=\$GOPATH/bin:\$PATH
go build -ldflags " -X \$PROJ/tests.BuildTag=\${BUILD_TAG} -X \$PROJ/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                    }
                }
                stage('windows_386') {
                    environment {
                        GOOS = "windows"
                        GOARCH = "386"
                        SUFFIX = ".exe"
                    }
                    steps {
                        sh """
export PATH=\$GOPATH/bin:\$PATH
go build -ldflags " -X \$PROJ/tests.BuildTag=\${BUILD_TAG} -X \$PROJ/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                    }
                }
                stage('windows_amd64') {
                    environment {
                        GOOS = "windows"
                        GOARCH = "amd64"
                        SUFFIX = ".exe"
                    }
                    steps {
                        sh """
export PATH=\$GOPATH/bin:\$PATH
go build -ldflags " -X \$PROJ/tests.BuildTag=\${BUILD_TAG} -X \$PROJ/tests.GitCommit=\$(git rev-parse HEAD)" -o build/\${GOOS}_\${GOARCH}/tech-db-forum\${SUFFIX}
cd build/\${GOOS}_\${GOARCH}
zip ../../target/dist/\${GOOS}_\${GOARCH}.zip tech-db-forum\${SUFFIX}
"""
                    }
                }
            }
        }
        stage('Prepare gp-pages') {
            steps {
                sh """
git branch -fD gh-pages || true
git branch -rd origin/gh-pages || true
ghp-import -n target/dist
"""
            }
        }
        stage('Publish gh-pages') {
            when {
                branch 'master'
            }
            steps {
                withCredentials([[$class: 'StringBinding', credentialsId: 'github_bozaro', variable: 'GITHUB_TOKEN']]) {
                    sh """
git push -qf https://\${GITHUB_TOKEN}@github.com/bozaro/tech-db-forum.git gh-pages
"""
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
                withCredentials([[$class: 'StringBinding', credentialsId: 'github_bozaro', variable: 'GITHUB_TOKEN']]) {
                    sh """
export PATH=\$GOPATH/bin:\$PATH
github-release info --tag ${params.TAG_NAME} || github-release release --tag ${params.TAG_NAME} --draft
for i in target/dist/*.zip; do
  github-release upload --tag ${params.TAG_NAME} --file \$i --name `basename \$i`
done
git tag ${params.TAG_NAME}
git push -qf https://\${GITHUB_TOKEN}@github.com/bozaro/tech-db-forum.git gh-pages ${params.TAG_NAME}
"""
                }
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: "target/dist/*.zip", fingerprint: true
        }
    }
}
