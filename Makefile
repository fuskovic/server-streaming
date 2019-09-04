name = file-server
image = fuskovic/server-streaming:1.0

image:

	@docker build -t ${image} .

stop :
	-docker stop ${name}
	-docker rm ${name}

container:

	@docker run --name ${name} -p 50051:50051 -d -t ${image}
	@docker exec -d ${name} server/run
	@echo "file server running in background on port 50051"

download:

	@go run client/client.go ${file}