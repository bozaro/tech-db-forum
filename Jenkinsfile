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
cd src/$goProject
go install ./vendor/github.com/bronze1man/yaml2json
go install ./vendor/github.com/aktau/github-release
"""
  }
  stage ('Build') {
    sh """#!/bin/bash -ex
export GOPATH="\$PWD"
export PATH="\$GOPATH/bin:\$PATH"

cd src/$goProject

# Build application
function go_build {
    GIT_COMMIT=`git rev-parse --verify HEAD`
    GOOS=\$1 GOARCH=\$2 go build -ldflags " -X github.com/bozaro/tech-db-forum/tests.BuildTag=\${BUILD_TAG} -X github.com/bozaro/tech-db-forum/tests.GitCommit=\${GIT_COMMIT}" -o build/\$1_\$2/tech-db-forum\$3

    pushd build/\$1_\$2
    zip ../../target/dist/\$1_\$2.zip tech-db-forum\$3
    popd
}

rm -fR target/dist/
mkdir -p target/dist/

go_build linux   amd64
go_build linux   386
go_build darwin  amd64
go_build windows amd64 .exe
go_build windows 386   .exe

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
