.PHONY: gen

gen:
	mkdir -p ./gen/go/achievements
	protoc --go_out=./gen/go/achievements --go_opt=paths=source_relative \
	--go-grpc_out=./gen/go/achievements --go-grpc_opt=paths=source_relative \
	./achievements.proto