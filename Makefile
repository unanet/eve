CI_COMMIT_BRANCH ?= local
CI_COMMIT_SHORT_SHA ?= 000001
CI_PROJECT_ID ?= 0
CI_PIPELINE_IID ?= 0
GOPATH ?= ${HOME}/go
MODCACHE ?= ${GOPATH}/pkg/mod

VERSION_MAJOR := 0
VERSION_MINOR := 1
VERSION_PATCH := 0
BUILD_NUMBER := ${CI_PIPELINE_IID}
PATCH_VERSION := ${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}
VERSION := ${PATCH_VERSION}.${BUILD_NUMBER}

DOCKER_UID = $(shell id -u)
DOCKER_GID = $(shell id -g)

CUR_DIR := $(shell pwd)

BUILD_IMAGE := unanet-docker.jfrog.io/golang
IMAGE_NAME := unanet-docker.jfrog.io/eve-api

LABEL_PREFIX := com.unanet
IMAGE_LABELS := \
	--label "${LABEL_PREFIX}.git_commit_sha=${CI_COMMIT_SHORT_SHA}" \
	--label "${LABEL_PREFIX}.gitlab_project_id=${CI_PROJECT_ID}" \
	--label "${LABEL_PREFIX}.build_number=${BUILD_NUMBER}" \
	--label "${LABEL_PREFIX}.version=${VERSION}"

docker-exec = docker run --rm \
	-e DOCKER_UID=${DOCKER_UID} \
	-e DOCKER_GID=${DOCKER_GID} \
	-v ${CUR_DIR}:/src \
	-v ${MODCACHE}:/go/pkg/mod \
	-v ${HOME}/.ssh/id_rsa:/home/unanet/.ssh/id_rsa \
	-w /src \
	${BUILD_IMAGE}

.PHONY: build dist

build:
	mkdir -p bin
	docker pull ${BUILD_IMAGE}
	$(docker-exec) go build ./...
	$(docker-exec) go test -short ./...
	$(docker-exec) go build -o ./bin/eve-api ./cmd/eve-api/main.go
	$(docker-exec) go build -o ./bin/eve-scheduler ./cmd/eve-scheduler/main.go

dist: build
	docker pull unanet-docker.jfrog.io/alpine-base
	docker build . -t ${IMAGE_NAME}:${PATCH_VERSION}
	docker push ${IMAGE_NAME}:${PATCH_VERSION}
	curl --fail -H "X-JFrog-Art-Api:${JFROG_API_KEY}" \
		-X PUT \
		https://unanet.jfrog.io/unanet/api/storage/docker-local/eve-api/${PATCH_VERSION}\?properties=version=${VERSION}%7Cgitlab-build-properties.project-id=${CI_PROJECT_ID}%7Cgitlab-build-properties.git-sha=${CI_COMMIT_SHORT_SHA}%7Cgitlab-build-properties.git-branch=${CI_COMMIT_BRANCH}
