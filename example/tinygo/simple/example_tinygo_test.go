package simple

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuperchain/xvm/compile"
	"github.com/xuperchain/xvm/debug"
	"github.com/xuperchain/xvm/exec"
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
	modulePath, err := compileLibrary("testdata/simple.wasm")
	if err != nil {
		log.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	code, err := exec.NewAOTCode(
		filepath.Join(cwd, modulePath),
		exec.NewMultiResolver(
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

	ret, err := ctx.Exec("_start", []int64{})

	if err != nil {
		log.Fatal(err)
	}
	_ = ret
	// Output: vim-go
}
