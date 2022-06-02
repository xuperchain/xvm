package rust

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuperchain/xvm/compile"
	"github.com/xuperchain/xvm/debug"
	"github.com/xuperchain/xvm/exec"
	"github.com/xuperchain/xvm/runtime/emscripten"
	"github.com/xuperchain/xvm/runtime/wasi"
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
	modulePath, err := compileLibrary("add.wasm")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(modulePath)
	fmt.Println(os.Getwd())
	code, err := exec.NewAOTCode(
		modulePath,
		exec.NewMultiResolver(
			emscripten.NewResolver(),
			wasi.NewResolver(),
		),
		&exec.CodeConfig{
			MemoryConfig: struct {
				MemoryGrow exec.MemoryGrowConfig
				Populate   exec.PopulateConfig
				HugePage   exec.HugePageConfig
			}{
				MemoryGrow: exec.MemoryGrowConfig{
					Enabled:    true,
					Initialize: 10, Maximium: 65535,
				},
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	ctx, err := code.NewContext(&exec.ContextConfig{GasLimit: exec.MaxGasLimit})
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.Release()
	debug.SetWriter(ctx, os.Stdout)

	ret, err := ctx.Exec("add", []int64{1, 2})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ret)
	// Output: 3
}
