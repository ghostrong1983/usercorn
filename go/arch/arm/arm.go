package arm

import (
	"github.com/lunixbochs/usercorn/go/models"

	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

func enterUsermode(u models.Usercorn) error {
	// place CPU into user mode
	modeSwitch := []byte{
		0x00, 0x00, 0x0f, 0xe1, // mrs r0, cpsr
		0x1f, 0x00, 0xc0, 0xe3, // bic r0, r0, $0x1f
		0x10, 0x00, 0x80, 0xe3, // orr r0, r0, $0x10
		0x00, 0xf0, 0x69, 0xe1, // msr spsr, r0
		0x0e, 0xf0, 0xb0, 0xe1, // movs pc, lr
	}
	lr, _ := u.RegRead(uc.ARM_REG_LR)
	sp, _ := u.RegRead(uc.ARM_REG_SP)
	mmap, err := u.Mmap(0, 0x1000)
	if err != nil {
		return err
	}
	end := mmap.Addr + uint64(len(modeSwitch))
	u.RegWrite(uc.ARM_REG_LR, end)
	u.MemWrite(mmap.Addr, modeSwitch)
	err = u.Start(mmap.Addr, end)
	u.MemUnmap(mmap.Addr, mmap.Size)
	u.RegWrite(uc.ARM_REG_LR, lr)
	u.RegWrite(uc.ARM_REG_SP, sp)
	return err
}