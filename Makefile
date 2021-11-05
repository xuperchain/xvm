build:
	make -C compile/wabt/
	mkdir bin ||true
	cp compile/wabt/build/wasm2c bin/

clean:
	@rm -rf bin
	@rm -rf compile/wabt/build