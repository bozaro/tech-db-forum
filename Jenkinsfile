goProject = "github.com/bozaro/tech-db-homework"

node ('linux') {
	stage ('Checkout') {
		checkout([
			$class: 'GitSCM',
			branches: scm.branches,
			doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
			extensions: scm.extensions + [
				[$class: 'CleanCheckout'],
			],
			userRemoteConfigs: scm.userRemoteConfigs
		])
	}
	stage ('Build') {
		sh """
# Generage swagger.json from swagger.yml
python -c 'import sys, yaml, json; json.dump(yaml.load(sys.stdin), sys.stdout, indent=4)' < swagger.yml > swagger.json
"""
	}
}
