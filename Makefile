clean:
	rm -f azure-nanny

build: clean
	CGO_ENABLED=0 go build ./cmd/azure-nanny

image: build
	go get github.com/openshift/imagebuilder/cmd/imagebuilder
	imagebuilder -f Dockerfile -t docker.io/mangirdas/azure-nanny:latest .

push:
	docker push docker.io/mangirdas/azure-nanny:latest

.PHONY: clean build push
