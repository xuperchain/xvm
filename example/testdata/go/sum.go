package main

import (
	"syscall/js"
)

func main() {
	js.Global().Get("a")
	js.Global().Get("b")
	value := 3
	js.Global().Set("c", js.ValueOf(value))
}
