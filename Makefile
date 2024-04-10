server:
	go run main.go

test:
	cd profile && go test -v
	# cd utils && go test -v

path:
	PATH="${PATH}:${HOME}/go/bin"