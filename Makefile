CI_COMMIT_BRANCH ?= local
CI_COMMIT_SHORT_SHA ?= 000001
CI_PROJECT_ID ?= 0
CI_PIPELINE_IID ?= 0

VERSION_MAJOR := 0
VERSION_MINOR := 1
VERSION_PATCH := 0
BUILD_NUMBER := ${CI_PIPELINE_IID}
PATCH_VERSION := ${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}
VERSION := ${PATCH_VERSION}.${BUILD_NUMBER}

DOCKER_IMAGE_NAME := unanet-docker.jfrog.io/eve-bot

build:
	docker build --ssh default . -t ${DOCKER_IMAGE_NAME}:${PATCH_VERSION}

dist: build
	docker push ${DOCKER_IMAGE_NAME}:${PATCH_VERSION}
	curl --fail -H "X-JFrog-Art-Api:${JFROG_API_KEY}" \
		-X PUT \
		https://unanet.jfrog.io/unanet/api/storage/docker/eve-bot/${PATCH_VERSION}\?properties=version=${VERSION}%7Cgitlab-build-properties.project-id=${CI_PROJECT_ID}%7Cgitlab-build-properties.git-sha=${CI_COMMIT_SHORT_SHA}%7Cgitlab-build-properties.git-branch=${CI_COMMIT_BRANCH}
