package rust

import (
	"fmt"

	"github.com/xuperchain/xvm/exec"
)

var resolver = exec.MapResolver(map[string]interface{}{
	"__wbindgen_placeholder__.__wbindgen_describe": func(ctx exec.Context) uint32 {
		return 0
	},
	"__wbindgen_placeholder__.__wbg_alert_531f0294104c16c5": func(ctx exec.Context, addr uint32, len uint32) uint32 {
		codec := exec.NewCodec(ctx)
		fmt.Println(codec.String(addr, len))
		return 0
	},
	"__wbindgen_externref_xform__.__wbindgen_externref_table_grow": func(ctx exec.Context, b uint32) uint32 {
		return 0
	},
	"__wbindgen_externref_xform__.__wbindgen_externref_table_set_null": func(ctx exec.Context, b uint32) {
	},
})

// NewResolver return exec.Resolver which resolves symbols needed by wasi environment
func NewResolver() exec.Resolver {
	return resolver
}
