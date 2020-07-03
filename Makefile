DOCKER_IMAGE_VERSION=v1
DOCKER_IMAGE=idirall22/gateway:${DOCKER_IMAGE_VERSION}

print:
	@echo DOCKER_IMAGE_VERSION = ${DOCKER_IMAGE_VERSION}
	@echo DOCKER_IMAGE = ${DOCKER_IMAGE}

gen:
	protoc -I proto/ proto/*.proto --go_out=plugins=grpc:pb/ \
	--grpc-gateway_out=:pb/ \
	--openapiv2_out=:open_api/

grpc:
	env SERVER_TYPE=grpc SERVER_PORT=:8080 go run cmd/server.go

rest:
	env SERVER_GRPC_ENDPOINT=0.0.0.0:8080 SERVER_PORT=:8081 SERVER_TYPE=rest go run cmd/server.go
	
client:	
	go run cmd/client.go

docker-build:
	docker build -t idirall22/gateway:v1 .

docker-push:
	docker push ${DOCKER_IMAGE}

k8s-clean:
	kubectl delete deploy --all

cert:
	cd cert && ./gen.cert && cd ..

.PHONY: cert k8s-clean docker-push docker-build client grpc rest gen