build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o mconfig-sdk-test main.go

docker: build
	docker build -t dockerhcy/mconfig-sdk-test:v0.1  .

push: docker
	docker push dockerhcy/mconfig-sdk-test:v0.1