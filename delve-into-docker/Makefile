.PHONY: build
build:
	docker build -t delve-into-docker-app .

.PHONY: run
run:
	docker run --rm \
		--publish 40000:40000 \
		--publish 1541:1541 \
		--security-opt=seccomp:unconfined \
		--name delve-into-docker \
		delve-into-docker-app
