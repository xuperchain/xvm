build:
	make -C compile/wabt/
	@ls bin>/dev/null ||mkdir bin
	@cp compile/wabt/build/wasm2c bin/

test: utest

utest:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...

clean:
	@rm -rf bin
	@rm -rf compile/wabt/build

