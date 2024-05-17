release:
	docker buildx build --platform linux/amd64,linux/arm64 --push -t afansv/vk-oauth-proxy:latest -t afansv/vk-oauth-proxy:$(IMAGE_TAG) .

