.PHONY: build run test docker-build docker-push deploy clean

# Variables
APP_NAME=go-mysql-app
DOCKER_REGISTRY=username
TAG=latest

build:
 go build -o bin/$(APP_NAME) cmd/main.go

run:
 go run cmd/main.go

test:
 go test -v ./...

docker-build:
 docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$(TAG) .

docker-push:
 docker push $(DOCKER_REGISTRY)/$(APP_NAME):$(TAG)

deploy-dev:
 kubectl apply -f k8s/
 kubectl rollout status deployment/go-app

deploy-prod:
 kubectl apply -f k8s/
 kubectl set image deployment/go-app app=$(DOCKER_REGISTRY)/$(APP_NAME):$(TAG)

clean:
 rm -rf bin/
 docker rmi $(DOCKER_REGISTRY)/$(APP_NAME):$(TAG)

logs:
 kubectl logs -f deployment/go-app

status:
 kubectl get all
 kubectl get pods