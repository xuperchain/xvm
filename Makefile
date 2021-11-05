build:
	make -C compile/wabt/
	mkdir bin ||true
	cp compile/wabt/build/wasm2c bin/

test: utest

utest:
	go test -coverprofile=coverage.txt -covermode=atomic ./...

clean:
	@rm -rf bin
	@rm -rf compile/wabt/build

