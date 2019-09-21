

build:
	GODEBUG="tls13=1" go build -o dogg3rz main.go 
	mkdir -p dist/
	mv dogg3rz dist/

install:
	mv dist/dogg3rz $$GOPATH/bin

