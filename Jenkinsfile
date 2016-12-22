/*
sudo apt install ghp-import
*/
node  ('linux') {
  stage ('Checkout') {
    checkout([
      $class: 'GitSCM',
      branches: scm.branches,
      doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
      extensions: scm.extensions + [
        [$class: 'CleanCheckout'],
        [$class: 'SubmoduleOption', disableSubmodules: false],
      ],
      userRemoteConfigs: scm.userRemoteConfigs
    ])
  }
  stage ('Prepare') {
    sh """
export GOPATH="\$PWD"
export PATH="\$GOPATH/bin:\$PATH"
go get github.com/bronze1man/yaml2json
"""
  }
  stage ('Build') {
    sh """
export GOPATH="\$PWD"
export PATH="\$GOPATH/bin:\$PATH"

# Generage swagger.json from swagger.yml
mkdir -p target
rm -fR target/doc/
cp -r swagger-ui/dist target/doc
yaml2json < swagger.yml > target/doc/swagger.json
sed -i 's/http:.*swagger.json/swagger.json/' target/doc/index.html
ghp-import -n target/doc
"""
  }
  if (env.BRANCH_NAME == 'master') {
    stage ('Publish') {
      withCredentials([usernamePassword(credentialsId: '88e000b8-d989-4f94-b919-1cc1352a5f96', passwordVariable: 'TOKEN', usernameVariable: 'LOGIN')]) {
        sh 'git push -qf https://${TOKEN}@github.com/bozaro/tech-db-forum.git gh-pages'
      }
    }
  }
}
