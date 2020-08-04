// +build linux

package chromedp

import (
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func allocateCmdOptions(cmd *exec.Cmd) {
	// In restricted environments (e.g. AWS Lambda), we are unable to set Pdeathsig
	if isPdeathsigRestricted() {
		return
	}

	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = new(syscall.SysProcAttr)
	}
	// When the parent process dies (Go), kill the child as well.
	cmd.SysProcAttr.Pdeathsig = syscall.SIGKILL
}

func isPdeathsigRestricted() bool {
	var sig int

	err := unix.Prctl(unix.PR_GET_PDEATHSIG, uintptr(unsafe.Pointer(&sig)), 0, 0, 0)
	if err != nil {
		return strings.Contains(err.Error(), "operation not permitted")
	}

	return false
}
