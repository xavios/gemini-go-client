run:
	rm -f gemini-client && go build && ./gemini-client

gen:
	go generate

mod:
	go mod tidy

live-test:
	CompileDaemon -build "go test -count=1"
