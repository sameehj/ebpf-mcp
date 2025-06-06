package ebpf

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/sys/unix"
)

func InspectSystemInfo() (map[string]interface{}, error) {
	kver, err := getKernelVersion()
	if err != nil {
		kver = fmt.Sprintf("unknown: %v", err)
	}

	support := map[string]interface{}{
		"go_arch":           runtime.GOARCH,
		"go_os":             runtime.GOOS,
		"kernel_version":    kver,
		"btf_enabled":       btfEnabled(),
		"sys_fs_bpf":        pathExists("/sys/fs/bpf"),
		"cgroup_v2_enabled": pathExists("/sys/fs/cgroup/cgroup.controllers"),
	}

	return support, nil
}

func getKernelVersion() (string, error) {
	var uname unix.Utsname
	if err := unix.Uname(&uname); err != nil {
		return "", err
	}
	return strings.TrimRight(string(uname.Release[:]), "\x00"), nil
}

func btfEnabled() bool {
	_, err := os.Stat("/sys/kernel/btf/vmlinux")
	return err == nil
}

func pathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}
