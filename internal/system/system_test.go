package system

import (
	"runtime"
	"testing"
)

func Test_SetCpuCore(t *testing.T) {
	SetCpuCore(0)
	SetCpuCore(-1)
	SetCpuCore(runtime.NumCPU())
}
