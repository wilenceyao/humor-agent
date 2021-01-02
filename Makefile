agent:
	go build -o humor-agent cmd/humor_agent.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/humoragent.proto
