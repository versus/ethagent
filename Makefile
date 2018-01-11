all: mac

mac:
	env GOOS=darwin GOARCH=amd64 go build  -a --ldflags="-s" -o ethagent-darwin-x64 .

linux:
	env GOOS=linux GOARCH=amd64 go build -a --ldflags="-s" -o ethagent-linux-x64 .

clean:
	rm ethagent-darwin-x64 || true
	rm ethagent-linux-x64 || true
