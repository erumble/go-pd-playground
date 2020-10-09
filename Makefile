service:
	go build -o bin/service github.com/erumble/go-api-boilerplate/cmd/service

clean:
	rm -rf bin/*
