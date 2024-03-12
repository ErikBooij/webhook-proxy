IMAGE_NAME=erikbooij/http-fan-out

.PHONY: build-image
build-image:
	$(eval VERSION := $(shell bash -c 'read -p "Version: " version; echo $$version'))
	echo "Building image $(IMAGE_NAME):$(VERSION)"
	docker build -t $(IMAGE_NAME):$(VERSION) .
	docker push $(IMAGE_NAME):$(VERSION)

.PHONY: build-image-latest
build-image-latest:
	$(eval VERSION := latest)
	echo "Building image $(IMAGE_NAME):$(VERSION)"
	docker build -t $(IMAGE_NAME):$(VERSION) .
	docker push $(IMAGE_NAME):$(VERSION)
