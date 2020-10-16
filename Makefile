service:
	go build -o bin/service github.com/erumble/go-pd-playground/cmd/service

clean:
	rm -rf bin/*
