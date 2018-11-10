build:
	docker build -t vision-app .

start:
	docker run --rm \
	--publish 8080:8080/tcp \
	--env GIN_MODE=release \
	--env AWS_ACCESS_KEY_ID=$$(aws configure get aws_access_key_id) \
	--env AWS_SECRET_ACCESS_KEY=$$(aws configure get aws_secret_access_key) \
	--name vision \
	vision-app
