PROJECT_ID = xxxx
IMAGE_NAME = xxxx
version ?= v0.0.1
VERSION = ${version}


GCR_URL = gcr.io/${PROJECT_ID}/${IMAGE_NAME}:${VERSION}

.PHONY: build
build:
	docker build -t ${GCR_URL} .


.PHONY: push
push:
	docker -- push ${GCR_URL}