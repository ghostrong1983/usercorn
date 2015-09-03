package x86_64

import (
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"

	"../../models"
	"../../syscalls"
)

var darwinSyscalls = map[int]string{
	1:   "exit",
	2:   "fork",
	3:   "read",
	4:   "write",
	5:   "open",
	6:   "close",
	7:   "wait4",
	9:   "link",
	10:  "unlink",
	73:  "munmap",
	197: "mmap",
	199: "lseek",
}

func DarwinSyscall(u models.Usercorn) {
	rax, _ := u.RegRead(uc.UC_X86_REG_RAX)
	ret, _ := u.Syscall(darwinSyscalls, int(rax-0x2000000), syscalls.RegArgs(u, AbiRegs))
	u.RegWrite(uc.UC_X86_REG_RAX, ret)
}

func DarwinInterrupt(u models.Usercorn, intno uint32) {
	if intno == 0x80 {
		DarwinSyscall(u)
	}
}

func init() {
	Arch.RegisterOS(&models.OS{Name: "darwin", Syscall: DarwinSyscall, Interrupt: DarwinInterrupt})
}
