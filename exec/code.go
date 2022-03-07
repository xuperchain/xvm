package exec

// #include "xvm.h"
// #include "stdlib.h"
// extern xvm_resolver_t make_resolver_t(void* env);
// #cgo LDFLAGS: -ldl
// #cgo linux CFLAGS: -DHUGEPAGE -D__USE_MISC
// #cgo linux CFLAGS: -DUSE_GNU
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/xuperchain/xvm/pointer"
)

type CodeConfig struct {
	MemoryConfig MemoryConfig
}
type MemoryGrowConfig struct {
	Enabled    bool
	Initialize int
	Maximium   int
}
type PopulateConfig struct {
	Enabled   bool
	Threshold int
}
type HugePageConfig struct {
	Enabled bool
}
type MemoryConfig struct {
	MemoryGrow MemoryGrowConfig
	Populate   PopulateConfig
	HugePage   HugePageConfig
}

var (
	DefaultCodeConfig = CodeConfig{
		MemoryConfig: MemoryConfig{
			MemoryGrow: MemoryGrowConfig{
				Enabled:    false,
				Initialize: 0,
				Maximium:   1024,
			},
			Populate: PopulateConfig{
				Enabled:   false,
				Threshold: 0,
			},
			HugePage: HugePageConfig{
				Enabled: false,
			},
		},
	}
)

// Code represents the wasm code object
type aotCode struct {
	code   *C.xvm_code_t
	bridge *resolverBridge
	// 因为cgo不能持有go的pointer，这个指针是一个指向bridge的token，最后需要Delete
	bridgePointer uintptr
}

// NewAOTCode instances a Code object from file path of native shared library
func NewAOTCode(module string, resolver Resolver, config *CodeConfig) (icode Code, err error) {
	code := new(aotCode)
	code.bridge = newResolverBridge(resolver)
	code.bridgePointer = pointer.Save(code.bridge)
	// xvm_new_code执行期间可能会抛出Trap，导致资源泄露
	// 如果CaptureTrap捕获了Trap则释放所有已经初始化的资源
	defer func() {
		if err != nil {
			code.Release()
			code = nil
		}
	}()
	defer CaptureTrap(&err)

	cpath := C.CString(module)
	defer C.free(unsafe.Pointer(cpath))
	resolvert := C.make_resolver_t(unsafe.Pointer(code.bridgePointer))
	code.code = C.xvm_new_code(cpath, resolvert)

	if code.code == nil {
		err = fmt.Errorf("open module %s error", module)
		return
	}

	cConfig := C.xvm_new_memory_config()
	cConfig.memory_grow_initialize = C.int(config.MemoryConfig.MemoryGrow.Initialize)
	cConfig.memory_grow_maximium = C.int(config.MemoryConfig.MemoryGrow.Maximium)

	if config.MemoryConfig.MemoryGrow.Enabled {
		cConfig.memory_grow_enabled = 1
	}

	if config.MemoryConfig.Populate.Enabled {
		cConfig.populate_enabled = 1
		cConfig.populate_threshold = C.int(config.MemoryConfig.Populate.Threshold)
	}

	if config.MemoryConfig.HugePage.Enabled {
		cConfig.huge_page_enabled = 1
	}

	ret := C.xvm_init_code(code.code, cConfig)
	if ret == 0 {
		err = fmt.Errorf("init module %s error", module)
		return
	}
	icode = code
	runtime.SetFinalizer(code, (*aotCode).Release)
	return
}

// Release releases resources hold by Code
func (c *aotCode) Release() {
	if c.code != nil {
		C.xvm_release_code(c.code)
	}
	if c.bridgePointer != 0 {
		pointer.Delete(c.bridgePointer)
	}
	*c = aotCode{}
	runtime.SetFinalizer(c, nil)
}
