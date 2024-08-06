VERSION  := $$(git describe --tags --always)
TARGET   := pod-spec-mutator
TEST     ?= ./...

default: test build

test:
	go test -v -run=$(RUN) $(TEST)

build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build \
		-a -tags netgo \
		-ldflags "-X main.Version=$(VERSION)" \
		-o bin/$(TARGET) .

publish:
	docker push mheers/$(TARGET):$(VERSION)
	docker tag mheers/$(TARGET):$(VERSION) mheers/$(TARGET):latest
	docker push mheers/$(TARGET):latest

image:
	docker build -t mheers/$(TARGET):$(VERSION) .

image-fast: build
	docker build -t mheers/$(TARGET):$(VERSION) -f Dockerfile.fast .

shell:
	docker run -it --entrypoint /bin/bash $(TARGET)

deploy-to-telepresence: build
	docker cp bin/pod-spec-mutator devcontainer-telepresence-1:/tmp/pod-spec-mutator

telepresence-get-certs:
	kubectl get secrets pod-spec-mutator-postgresoperator-tls -o json | jq -r '.data."tls.crt"' | base64 -d > /tmp/tls.crt
	kubectl get secrets pod-spec-mutator-postgresoperator-tls -o json | jq -r '.data."tls.key"' | base64 -d > /tmp/tls.key
	docker exec -it devcontainer-telepresence-1 mkdir -p /tmp/k8s-webhook-server/serving-certs
	docker cp /tmp/tls.crt devcontainer-telepresence-1:/tmp/k8s-webhook-server/serving-certs/tls.crt
	docker cp /tmp/tls.key devcontainer-telepresence-1:/tmp/k8s-webhook-server/serving-certs/tls.key

exec-telepresence:
	docker exec -it devcontainer-telepresence-1 bash
