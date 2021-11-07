TODAY=$(shell date +%y-%m-%d)

serve:
	go run server.go

builder:
	docker build \
 		-t builder docker/build/


bin/cms:
	docker run -v ${PWD}../../..:/src -w /src -it builder \
	env GOARCH=amd64 go build \
	-ldflags "-X main.version=${VERSION}" \
	-v -o cmd/serverdb/bin/cms cmd/serverdb/server.go


clean:
	@rm -f bin/*


docker:
	docker volume create cms
	docker network create web


run:
	docker run \
	  --mount source=cms,target=/app/vol \
	  --network web \
	  --env-file deploy.env \
	  -p 80:80\
	  --name cms -it --rm cms:${TODAY} bash

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
	   --name cms -d cms:${TODAY} ./cms


stop:
	docker stop cms


rm:
	docker rm cms


logs:
	docker logs -f cms

