build:
	make -C compile/wabt/
	@ls bin 2>/dev/null 1>/dev/null || mkdir bin 
	cp compile/wabt/build/wasm2c bin/

test: utest spec 

utest:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
	
spec:
	bash -c "cd spectest && go run main.go ./core && cd .."

clean:
	@rm -rf bin
	@rm -rf compile/wabt/build

