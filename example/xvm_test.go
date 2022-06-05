package example

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"

	"github.com/xuperchain/xvm/compile"
	"github.com/xuperchain/xvm/debug"
	"github.com/xuperchain/xvm/exec"
	gowasm "github.com/xuperchain/xvm/runtime/go"
)

func compileLibrary(wasmpath string) (string, error) {
	tmpdir, err := ioutil.TempDir("", "xvm-exec-test")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpdir)
	cfg := &compile.Config{
		Wasm2cPath: "wasm2c",
		OptLevel:   0,
	}
	libpath := replaceExt(wasmpath, ".so")
	err = compile.CompileNativeLibrary(cfg, libpath, wasmpath)
	if err != nil {
		return "", err
	}
	return libpath, nil
}

func replaceExt(name, ext string) string {
	dir, file := filepath.Split(name)
	idx := strings.Index(file, ".")
	if idx == -1 {
		file = file + ext
	} else {
		file = file[:idx] + ext
	}
	return filepath.Join(dir, file)
}

func ExampleAdd() {
	modulePath, err := compileLibrary("testdata/add.wasm")
	if err != nil {
		log.Fatal(err)
	}
	code, err := exec.NewAOTCode(modulePath, nil, &exec.DefaultCodeConfig)
	if err != nil {
		log.Fatal(err)
	}
	ctx, err := code.NewContext(&exec.ContextConfig{GasLimit: exec.MaxGasLimit})
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Release()
	ret, err := ctx.Exec("_add", []int64{1, 2})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ret)
	// Output: 3
}

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
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(path)
	modulePath, err := compileLibrary(path)
	if err != nil {
		log.Fatal(err)
	}
	resolver := gowasm.NewResolver()
	code, err := exec.NewAOTCode(
		filepath.Join(cwd, modulePath),
		resolver,
		&exec.DefaultCodeConfig,
	)
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
