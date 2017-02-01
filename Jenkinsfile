/*
sudo apt install ghp-import
*/
goProject = "github.com/bozaro/tech-db-forum"

properties([parameters([string(name: 'TAG_NAME', defaultValue: '')])])
if (params.TAG_NAME != "") {
  echo "Build tag: ${params.TAG_NAME}"
}

node  ('linux') {
  stage ('Checkout') {
    checkout([
      $class: 'GitSCM',
      branches: scm.branches,
      doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
      extensions: scm.extensions + [
        [$class: 'CleanCheckout'],
        [$class: 'RelativeTargetDirectory', relativeTargetDir: "src/$goProject"],
        [$class: 'SubmoduleOption', disableSubmodules: false, recursiveSubmodules: false],
      ],
      userRemoteConfigs: scm.userRemoteConfigs
    ])
  }
  stage ('Prepare') {
    sh """
export GOPATH="\$PWD"
export PATH="\$GOPATH/bin:\$PATH"
go get github.com/bronze1man/yaml2json
go get github.com/aktau/github-release
"""
  }
  stage ('Build') {
    sh """#!/bin/bash -ex
export GOPATH="\$PWD"
export PATH="\$GOPATH/bin:\$PATH"

cd src/$goProject

# Build application
go build
GOOS=linux	GOARCH=amd64	; go build -o build/\${GOOS}_\${GOARCH}/tech-db-forum
GOOS=linux	GOARCH=386	; go build -o build/\${GOOS}_\${GOARCH}/tech-db-forum
GOOS=darwin	GOARCH=amd64	; go build -o build/\${GOOS}_\${GOARCH}/tech-db-forum
GOOS=darwin	GOARCH=386	; go build -o build/\${GOOS}_\${GOARCH}/tech-db-forum
GOOS=windows	GOARCH=amd64	; go build -o build/\${GOOS}_\${GOARCH}/tech-db-forum.exe
GOOS=windows	GOARCH=386	; go build -o build/\${GOOS}_\${GOARCH}/tech-db-forum.exe

rm -fR target/dist/
mkdir -p target/dist/
for i in build/*/; do
  pushd \$i
  zip ../../target/dist/`basename \$i`.zip *
  popd
done
git branch -fD gh-pages || true
git branch -rd origin/gh-pages || true
ghp-import -n target/dist
"""
    archive "src/$goProject/target/dist/*.zip"
  }
  if (env.BRANCH_NAME == 'master') {
    stage ('Publish') {
      withCredentials([usernamePassword(credentialsId: '88e000b8-d989-4f94-b919-1cc1352a5f96', passwordVariable: 'TOKEN', usernameVariable: 'LOGIN')]) {
        sh """
cd src/$goProject
git push -qf https://\${TOKEN}@github.com/bozaro/tech-db-forum.git gh-pages
"""
      }
    }
  }
  if (params.TAG_NAME != "") {
    stage ("Publish: github") {
      withEnv([
        "TAG_NAME=${params.TAG_NAME}",
        "GITHUB_USER=bozaro",
        "GITHUB_REPO=tech-db-forum",
      ]) {
        withCredentials([[$class: 'StringBinding', credentialsId: '49bf22be-f4d4-4a75-855a-b0e56e357f1c', variable: 'GITHUB_TOKEN']]) {
          sh """
export GOPATH="\$PWD"
export PATH="\$GOPATH/bin:\$PATH"

github-release info --tag \$TAG_NAME || github-release release --tag \$TAG_NAME --draft
for i in src/$goProject/target/dist/*.zip; do
  github-release upload --tag \$TAG_NAME --file \$i --name `basename \$i`
done
"""
        }
      }
    }
  }
}
