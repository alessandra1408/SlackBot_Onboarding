include scripts/help.mk

.PHONY: run-local git-config version clean install lint audit env env-ip env-stop test cover build image tag push release run run-docker remove-docker
.DEFAULT_GOAL := help

GITLAB_GROUP   	= consulting-service
REGISTRY      	= gcr.io
REGISTRY_GROUP 	= sebrae-ddm-prod-7518

BUILD         	= $(shell git rev-parse --short HEAD)
DATE          	= $(shell date -uIseconds)
VERSION  	  	= $(shell git describe --always --tags)
NAME           	= $(shell basename $(CURDIR))
IMAGE          	= $(REGISTRY)/$(REGISTRY_GROUP)/$(NAME):$(BUILD)

POSTGRES_NAME   = postgres_$(NAME)_$(BUILD)
PUB_SUB_EMULATOR = pub_sub_emulator_$(NAME)_$(BUILD)
NETWORK_NAME  = network_$(NAME)_$(BUILD)

git-config:
	git config --replace-all core.hooksPath .githooks

check-env-%:
	@ if [ "${${*}}" = ""  ]; then \
		echo "Variable '$*' not set"; \
		exit 1; \
	fi

version: ##@other Check version.
	@echo $(VERSION)

clean: ##@dev Remove folder vendor, public and coverage.
	rm -rf vendor public coverage

install: clean ##@dev Download dependencies via go mod.
	GO111MODULE=on go mod download
	GO111MODULE=on go mod vendor

audit: clean ##@check Run vulnerability check in Go dependencies.
	DOCKER_BUILDKIT=1 docker build --progress=plain --target=audit --file=./build/Dockerfile .

lint: clean ##@check Run lint on docker.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--target=lint \
		--file=./build/Dockerfile .

env: ##@environment Create network and run postgres container.
	POSTGRES_NAME=${POSTGRES_NAME} \
	PUB_SUB_EMULATOR=${PUB_SUB_EMULATOR} \
	NETWORK_NAME=${NETWORK_NAME} \
	CURRENT_UID=$(id -u):$(id -g) docker-compose -f ./test/docker-compose.yml up -d

env-ip: ##@environment Return local Postgres IP (from Docker container)
	@echo POSTGRES $$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${POSTGRES_NAME})
	@echo PUB_SUB_EMULATOR $$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${PUB_SUB_EMULATOR})

env-stop: ##@environment Remove postgres container and remove network.
	POSTGRES_NAME=${POSTGRES_NAME} PUB_SUB_EMULATOR=${PUB_SUB_EMULATOR} NETWORK_NAME=${NETWORK_NAME} docker-compose -f ./test/docker-compose.yml kill
	POSTGRES_NAME=${POSTGRES_NAME} PUB_SUB_EMULATOR=${PUB_SUB_EMULATOR} NETWORK_NAME=${NETWORK_NAME} docker-compose -f ./test/docker-compose.yml rm -vf
	docker network rm $(NETWORK_NAME)

test: clean ##@check Run tests and coverage.
	docker build --progress=plain \
		--network $(NETWORK_NAME) \
		--tag $(IMAGE) \
		--build-arg PUB_SUB_EMULATOR_HOST=$$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(PUB_SUB_EMULATOR)):8085 \
		--build-arg POSTGRES_URL=postgres://postgres:admin@$(POSTGRES_NAME):5432/sebrae?sslmode=disable \
		--target=test \
		--file=./build/Dockerfile .

	-mkdir coverage
	docker create --name $(NAME)-$(BUILD) $(IMAGE)
	docker cp $(NAME)-$(BUILD):/index.html ./coverage/.
	docker rm -vf $(NAME)-$(BUILD)

build: clean ##@build Build image.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--tag $(IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--target=build \
		--file=./build/Dockerfile .

image: check-env-VERSION clean ##@build Create release docker image.
	DOCKER_BUILDKIT=1 \
	docker build --progress=plain \
		--tag $(IMAGE) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD=$(BUILD) \
		--build-arg DATE=$(DATE) \
		--target=image \
		--file=./build/Dockerfile .

tag: check-env-VERSION ##@build Add docker tag.
	docker tag $(IMAGE) \
		$(REGISTRY)/$(REGISTRY_GROUP)/$(NAME):$(VERSION)

push: check-env-VERSION ##@build Push docker image to registry.
	docker push $(REGISTRY)/$(REGISTRY_GROUP)/$(NAME):$(VERSION)

release: check-env-TAG ##@build Create and push git tag.
	git tag -a $(TAG) -m "Generated release "$(TAG)
	git push origin $(TAG)


run: ##@dev Run locally.
	go run cmd/app/main.go

run-local: ##@dev Run locally.
	ENV_LOGGER_LEVEL=debug \
	ENV_WORKERS=1 \
	ENV_GCP_POSTGRES_HOST=localhost \
	ENV_GCP_POSTGRES_PORT=15432 \
	ENV_GCP_POSTGRES_USER=sebrae_ingestion \
	ENV_GCP_POSTGRES_DATABASE=datalake \
	ENV_GCP_CRD=/home/rudsonsouza/.config/gcloud/application_default_credentials.json \
	ENV_GCP_PROJECT_ID=sebrae-ddm-qa-3746 \
	ENV_GCP_PUBSUB_SUBSCRIPTION_ID=marketing-cloud-error-handling-sub-qa \
	ENV_GCP_PUBSUB_TOPIC_ID=marketing-cloud-error-handling-qa \
	ENV_MARKETING_CLOUD_URL_BASE=mcdgv4q4mfzp49dnzqs2hr-gjqm0 \
	go run cmd/app/main.go

run-docker: check-env-POSTGRES_URL ##@docker Run docker container.
	docker run --rm \
		--name $(NAME) \
		-e LOGGER_LEVEL=debug \
		-e POSTGRES_URL=$(POSTGRES_URL) \
		-p 5001:8080 \
		$(IMAGE)

remove-docker: ##@docker Remove docker container.
	-docker rm -vf $(NAME)
