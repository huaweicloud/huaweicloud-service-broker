.PHONY: all build huaweicloud-service-broker docker clean

all:build

build:huaweicloud-service-broker

huaweicloud-service-broker:
	go build

docker:huaweicloud-service-broker
	docker build . -t quay.io/huaweicloud/huaweicloud-service-broker:latest

clean:
	rm -rf ./huaweicloud-service-broker
