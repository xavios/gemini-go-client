run:
	rm -f gemini-client && go build && ./gemini-client

gen:
	go generate

mod:
	go mod tidy

live-test:
	CompileDaemon -build "go clean -cache && go test"
