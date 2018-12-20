.PHONY: xray

xray:
	docker run --rm \
		--env AWS_ACCESS_KEY_ID=$$(aws configure get aws_access_key_id) \
		--env AWS_SECRET_ACCESS_KEY=$$(aws configure get aws_secret_access_key) \
		--env AWS_REGION=eu-west-2 \
		--name xray-daemon \
		--publish 2000:2000/udp \
		amazon/aws-xray-daemon -o
