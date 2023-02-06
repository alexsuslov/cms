TODAY=$(shell date +%y-%m-%d)

serve:
	go run server.go

builder:
	docker build \
 		-t builder docker/build/


bin/wiki:
	docker run -v ${PWD}../../..:/src -w /src -it builder \
	env GOARCH=amd64 go build \
	-ldflags "-X main.version=${VERSION}" \
	-v -o cmd/wiki/bin/wiki cmd/serverdb/server.go


clean:
	@rm -f bin/*


docker:
	docker volume create wiki
	docker network create web


run:
	docker run \
	  --mount source=cms,target=/app/vol \
	  --network web \
	  --env-file deploy.env \
	  -p 80:80\
	  --name wiki -it --rm wiki:${TODAY} bash

copy:
	docker cp templates cms:/app/vol/
	docker cp -a static cms:/app/vol/
	docker cp config.yml cms:/app/vol/


start:
	docker run \
	  --mount source=cms,target=/app/vol \
	  --network web \
	  --env-file deploy.env \
	  -p 80:80\
	   --name wiki -d wiki:${TODAY} ./wiki


stop:
	docker stop wiki


rm:
	docker rm wiki


logs:
	docker logs -f wiki

