package example

import (
	"fmt"
	"log"
	"os"
	osexec "os/exec"
	"path/filepath"

	"github.com/xuperchain/xvm/debug"
	"github.com/xuperchain/xvm/exec"
	gowasm "github.com/xuperchain/xvm/runtime/go"
)

func ExampleGo() {
	path := "testdata/go/sum"

	cmd := osexec.Command("go", "build", "sum.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{"GOOS=js", "GOARCH=wasm"}
	cmd.Env = append(os.Environ(), cmd.Env...)
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Dir = filepath.Join(cwd, "testdata/go")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(path)
	modulePath, err := compileLibrary(path)
	if err != nil {
		log.Fatal(err)
	}
	resolver := gowasm.NewResolver()
	code, err := exec.NewAOTCode(modulePath, resolver)
	if err != nil {
		log.Fatal(err)
	}
	ctx, err := code.NewContext(&exec.ContextConfig{GasLimit: exec.MaxGasLimit})
	if err != nil {
		log.Fatal(err)
	}
	debug.SetWriter(ctx, os.Stderr)
	rt := gowasm.RegisterRuntime(ctx)
	rt.Global().Register("a", 1)
	rt.Global().Register("b", 2)

	defer ctx.Release()
	ret, err := ctx.Exec("run", []int64{1, 2})
	if err != nil {
		log.Fatal(err)
	}
	if ret != 0 {
		log.Fatalf("expert %d, got %d\n", 0, ret)
	}
	c, exist := rt.Global().GetProperty("c")
	if !exist {
		log.Fatal(exist)
	}
	fmt.Println(c)
	// 	Output: 3
}
