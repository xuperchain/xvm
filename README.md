# XuperVM

XuperVM 是 [XuperChain](https://xuper.baidu.com/n/ps/opensource) 旗下的高性能 WebAssembly 虚拟机，采用 AOT 编译技术，支持资源审计与资源限制，原生支持 Golang 运行时

## 快速开始
### 构建
```shell script
    git https://github.com/xuperchain/xvm.git
    cd xvm 
    make
    # 设置环境变量
    export PATH=`pwd`/bin:${PATH}
```

> 也可以将环境变量修改写入到 bashrc(bashprofile)，以便每次启动时新终端时自动设置环境变量


### 运行第一个WASM 程序

以 [add.wasm](example/testdata/add.wasm) 为例

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuperchain/xvm/compile"
	"github.com/xuperchain/xvm/exec"
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

func main() {
	modulePath, err := compileLibrary("example/testdata/add.wasm")
	if err != nil {
		log.Fatal(err)
	}
	resolver := emscripten.NewResolver()
	code, err := exec.NewAOTCode(modulePath, resolver)
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
```


更多的例子可以在 [example](example) 目录里查看


## 参与贡献
我们欢迎任何形式的贡献(新的功能/文档/Issue)，如果你想要为XuperVM 贡献代码，你需要
1. 查看并签署[贡献者协议](https://cla-assistant.io/xuperchain/xvm)
2. 查看 [代码提交指南](https://xuper.baidu.com/n/xuperdoc/contribution/pull_requests.html)
3. 本地编辑文件
4. 提交完整的 Pull Request

## 联系我们
商务合作，请Email：xchain-help@baidu.com, 来源请注明Github。
如果你对XuperChain开源技术及应用感兴趣，欢迎添加“百度超级链·小助手“微信，回复“技术论坛进群”，加入“百度超级链开发者社区”，与百度资深工程师深度交流!微信二维码如下:

![微信二维码](https://github.com/ToWorld/xuperchain-image/blob/master/baidu-image-xuperchain.png)