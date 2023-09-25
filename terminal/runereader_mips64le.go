//go:build mips64le && linux
// +build mips64le,linux

package terminal

// Used syscall numbers from https://github.com/golang/go/blob/master/src/syscall/ztypes_linux_mips64le.go
const ioctlReadTermios = 0x540d  // syscall.TCGETS
const ioctlWriteTermios = 0x540e // syscall.TCSETS
