package exec

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/xuperchain/xvm/compile"
)

func withCode(t testing.TB, watCode string, r Resolver, config CodeConfig, f func(code Code)) {
	tmpdir, err := ioutil.TempDir("", "xvm-exec-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	wasmpath := filepath.Join(tmpdir, "wasm.wasm")
	libpath := filepath.Join(tmpdir, "wasm.so")
	cfg := &compile.Config{
		Wasm2cPath:   "../compile/wabt/build/wasm2c",
		Wat2wasmPath: "../compile/wabt/build/wat2wasm",
		OptLevel:     0,
	}

	err = compile.CompileWatSource(cfg, wasmpath, watCode)
	if err != nil {
		t.Fatal(err)
	}

	err = compile.CompileNativeLibrary(cfg, libpath, wasmpath)
	if err != nil {
		t.Fatal(err)
	}
	code, err := NewAOTCode(libpath, r, &config)
	if err != nil {
		t.Fatal(err)
	}
	f(code)
	code.Release()
}

func TestNewCode(t *testing.T) {
	withCode(t, "testdata/sum.wat", nil, DefaultCodeConfig, func(code Code) {
	})
}

func BenchmarkNewCode(b *testing.B) {
	b.Run("populate", func(b *testing.B) {
		config := CodeConfig{}
		config.MemoryConfig.Populate.Enabled = true
		withCode(b, "testdata/add.wat", nil, config, func(code Code) {
			for i := 0; i < b.N; i++ {
				ctx, err := code.NewContext(&ContextConfig{GasLimit: MaxGasLimit})
				if err != nil {
					b.Fatal(err)
				}
				ctx.Exec("_add", []int64{1, 2})
				ctx.Release()
			}
		})
	})

	b.Run("BaseLine", func(b *testing.B) {
		config := CodeConfig{}
		config.MemoryConfig.Populate.Enabled = false
		withCode(b, "testdata/add.wat", nil, config, func(code Code) {
			for i := 0; i < b.N; i++ {
				ctx, err := code.NewContext(&ContextConfig{GasLimit: MaxGasLimit})
				if err != nil {
					b.Fatal(err)
				}
				ctx.Exec("_add", []int64{1, 2})
				ctx.Release()
			}
		})
	})
}
