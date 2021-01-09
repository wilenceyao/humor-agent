agent:
	go build -o humor-agent cmd/humor_agent.go

protoc:
	protoc --go_out=. --go_opt=paths=source_relative api/humor/api.proto