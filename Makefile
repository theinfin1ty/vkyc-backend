# Makefile

# Define the name of your Docker image
IMAGE_NAME := vkyc_backend

# Define the name for your container
CONTAINER_NAME := vkyc

# Build the Docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run the Docker container
run:
	docker run --name $(CONTAINER_NAME) --net=host --rm $(IMAGE_NAME)

# Stop and remove the Docker container
clean:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

# Full cycle: build, run, and clean
build-run: build run

.PHONY: build run clean all