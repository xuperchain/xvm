package wasi

import (
	"unsafe"

	"github.com/xuperchain/xvm/debug"
	"github.com/xuperchain/xvm/exec"
)

func errno(n int32) uint32 {
	return *(*uint32)(unsafe.Pointer(&n))
}

var resolver = exec.MapResolver(map[string]interface{}{
	"wasi_snapshot_preview1.fd_prestat_get": func(ctx exec.Context, x, y uint32) uint32 {
		return 8
	},
	"wasi_snapshot_preview1.fd_fdstat_get": func(ctx exec.Context, x, y uint32) uint32 {
		return 8
	},
	"wasi_snapshot_preview1.fd_prestat_dir_name": func(ctx exec.Context, x, y, z uint32) uint32 {
		return 8
	},
	"wasi_snapshot_preview1.fd_close": func(ctx exec.Context, x uint32) uint32 {
		return 8
	},
	"wasi_snapshot_preview1.fd_seek": func(ctx exec.Context, x, y, z, w uint32) uint32 {
		return 8
	},
	"wasi_snapshot_preview1.fd_write": func(ctx exec.Context, fd, y, iovcnt, w uint32) uint32 {
		codec := exec.NewCodec(ctx)
		// only stdout and stderr is supported
		if fd != 1 && fd != 2 {
			return errno(-9)
		}
		total := uint32(0)
		for i := uint32(0); i < iovcnt; i++ {
			base := codec.Uint32(y + 8*i)
			length := codec.Uint32(y + 8*i + 4)
			buf := codec.Bytes(base, length)
			debug.Write(ctx, buf)
			total += length
		}
		return total
	},
	"wasi_snapshot_preview1.environ_sizes_get": func(ctx exec.Context, x, y uint32) uint32 {
		return 0
	},
	"wasi_snapshot_preview1.environ_get": func(ctx exec.Context, x, y uint32) uint32 {
		return 0
	},
	"wasi_snapshot_preview1.args_sizes_get": func(ctx exec.Context, x, y uint32) uint32 {
		return 0
	},
	"wasi_snapshot_preview1.args_get": func(ctx exec.Context, x, y uint32) uint32 {
		return 0
	},
	"wasi_snapshot_preview1.fd_read": func(ctx exec.Context, a, b, c, d uint32) uint32 {
		return 0
	},

	"wasi_snapshot_preview1.proc_exit": func(ctx exec.Context, x uint32) uint32 {
		exec.Throw(exec.NewTrap("exit"))
		return 0
	},
	"wasi_snapshot_preview1.clock_time_get": func(exec.Context, uint32, uint64, uint32) uint32 {
		return 0
	},
})

// NewResolver return exec.Resolver which resolves symbols needed by wasi environment
func NewResolver() exec.Resolver {
	return resolver
}
