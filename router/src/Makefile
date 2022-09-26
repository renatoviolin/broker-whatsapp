mongo-up: #create mongo container
	docker run --name mongodb -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=renato -e MONGO_INITDB_ROOT_PASSWORD=renato mongo
# docker run --name mongodb -v /tmp/mongo-db:/data/db -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=renato -e MONGO_INITDB_ROOT_PASSWORD=renato mongo

mongo-down: #destroy mongo container
	docker rm -f mongodb

build: #destroy mongo container
	GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o main cmd/main.go

run-local: # run localhost
	go run cmd/main.go