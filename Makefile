run-user-server:
	go run user-service/main.go

run-user-client:
	go run user-service/client/main.go

run-notification-server:
	go run notification-service/main.go

run-notification-client:
	go run notification-service/client/main.go

rpc-user:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user.proto
rpc-notification:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/notification.proto